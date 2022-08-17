package opensdk

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/junhwong/go-logs"
	log "github.com/junhwong/go-logs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	ErrNotExecuteYet = errors.New("not execute yet.")
)

type HttpClient interface {
	Do(request *http.Request, tlsTwowayAuthentication bool) (*http.Response, error)
}

// Executor 用于执行请求的相关上下文。
// 每个 Executor 都可以被多次执行(注意：接口业务是否允许)。
type Executor interface {
	// Execute 执行请求并返回结果。
	// verifySign 表示是否需要同步验证签名，默认 `false` 不验证。
	Execute(verifySign ...bool) Results
	SetNotifyURL(url string, filedName ...string) Executor //设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	UseTwowayAuthentication(b bool) Executor               // 是否需要双向认证
	UseXML(b bool) Executor                                //是否使用xml作为接口数据交换格式
	ResultValidator(f func(Params) (ok bool, code string, msg string, subcode string, submsg string)) Executor
	Set(filed string, value interface{}) Executor
	//WithDecoder(decoder func(data []byte, dataFormat string, out *Params) (err error))
}

type DefaultExecutor struct {
	Params
	Client           Client
	BuildRequest     func(*DefaultExecutor) (req *http.Request, err error)
	successValidator func(Params) (ok bool, code string, msg string, subcode string, submsg string)
	Decoder          func(data []byte, dataFormat string, out *Params) (err error)
	TLS              bool
	DataFormat       string
	Err              error
	APIURL           string
	HTTPMethod       string
}

func (e *DefaultExecutor) ResultValidator(f func(Params) (ok bool, code string, msg string, subcode string, submsg string)) Executor {
	e.successValidator = f
	return e
}
func (e *DefaultExecutor) UseTwowayAuthentication(b bool) Executor {
	e.TLS = b
	return e
}
func (e *DefaultExecutor) Set(filed string, value interface{}) Executor {
	e.Params[filed] = value
	return e
}
func (e *DefaultExecutor) UseXML(b bool) Executor {
	if b {
		e.DataFormat = "xml"
	} else {
		e.DataFormat = "default"
	}
	return e
}

func (e *DefaultExecutor) Execute(verifySign ...bool) (res Results) {
	r := DefaultResults{Params: Params{}, Err: e.Err}
	res = &r
	if r.Err != nil {
		return
	}
	defer func() {
		if x := recover(); x != nil {
			if err, ok := x.(error); ok {
				e.Err = err
				log.Print(err)
			} else {
				log.Print(x)
			}
		}
		r.Err = e.Err
	}()
	// println("here10")
	// TODO: 计时
	hc, err := e.Client.HttpClient(e.TLS)
	if err != nil {
		// println("here12", err.Error(), hc)
		e.Err = err
		return
	}
	// println("here11")
	request, err := e.BuildRequest(e)
	if err != nil {
		e.Err = err
		return
	}
	// 不使用这个会产生 EOF 错误 !! see: https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi/23963271#23963271
	request.Close = true
	// println("here1")
	response, err := hc.Do(request)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	//
	// println("here")
	if err != nil {
		log.Printf("执行请求错误：%+v", err)
		e.Err = err
		return
	}
	r.Params["response.StatusCode"] = response.StatusCode
	r.Params["response.Header"] = response.Header
	var tr io.Reader = response.Body
	if response.Header != nil {
		arr := strings.Split(response.Header.Get("Content-Type"), ";")
		for _, s := range arr {
			if arr := strings.SplitN(s, "=", 2); len(arr) == 2 &&
				strings.EqualFold(strings.ToLower(arr[0]), "charset") &&
				strings.EqualFold(strings.ToLower(arr[1]), "gbk") {
				tr = transform.NewReader(tr, simplifiedchinese.GBK.NewDecoder())
				break
			}
		}
	}
	r.Data, r.Err = ioutil.ReadAll(tr)
	if r.Err != nil {
		log.Print(r.Err)
		return
	}
	// fmt.Println(string(r.Data))
	logs.Prefix("go-opensdk").Debug("response: ", string(r.Data))
	// fmt.Println(r)

	if e.Decoder == nil {
		e.Decoder = DefaultDecoder
	}
	r.Err = e.Decoder(r.Data, e.DataFormat, &r.Params)
	if r.Err != nil {
		if r.Err == io.EOF {
			e.Err = nil
		} else {
			return
		}
	}
	if len(verifySign) > 0 && verifySign[0] {
		// TODO: 验证签名
	}
	if e.successValidator != nil {
		r.ResultSuccess, r.ResultCode, r.ResultMsg, r.ResultSubCode, r.ResultSubMsg = e.successValidator(r.Params)
	}

	return
}

//SetNotifyURL 设置接口服务器主动通知调用服务器里指定的页面http/https路径。
func (e *DefaultExecutor) SetNotifyURL(url string, filedName ...string) Executor {
	filed := "notify_url"
	if len(filedName) > 0 && filedName[0] != "" {
		filed = filedName[0]
	}
	e.Params[filed] = url
	return e
}

func DefaultDecoder(data []byte, dataFormat string, out *Params) (err error) {

	switch dataFormat {
	case "xml":
		err = xml.Unmarshal(data, out)
	default:
		err = json.Unmarshal(data, out)
	}

	return
}
