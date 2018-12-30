package opensdk

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrNotExecuteYet = errors.New("not execute yet.")
)

// Executor 用于执行请求的相关上下文。
// 每个 Executor 都可以被多次执行(注意：接口业务是否允许)。
type Executor interface {
	// Execute 执行请求并返回结果。
	// verifySign 表示是否需要同步验证签名，默认 `false` 不验证。
	Execute(verifySign ...bool) Results
	SetNotifyURL(url string, filedName ...string) Executor //设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	UseTLS(b bool) Executor                                // 是否需要双向认证
	UseXML(b bool) Executor                                //是否使用xml作为接口数据交换格式
	ResultValidator(f func(Params) (ok bool, code string, msg string, subcode string, submsg string)) Executor
	Set(filed string, value interface{}) Executor
}

type DefaultExecutor struct {
	Params
	Client           Client
	Request          func(*DefaultExecutor) (response *http.Response, requestLog string, err error)
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
func (e *DefaultExecutor) UseTLS(b bool) Executor {
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
	// defer func() {
	// 	if x := recover(); x != nil {
	// 		if err, ok := x.(error); ok {
	// 			r.Err = err
	// 		} else {
	// 			log.Print(x)
	// 		}
	// 	}
	// 	// 日志
	// }()

	// TODO: 计时
	// var log string
	response, _, err := e.Request(e)
	if err != nil {
		panic(err)
	}
	r.Data, r.Err = ioutil.ReadAll(response.Body)
	if response.Body != nil {
		response.Body.Close()
	}
	if r.Err == io.EOF {
		e.Err = nil
	}

	if e.Decoder == nil {
		e.Decoder = DefaultDecoder
	}
	err = e.Decoder(r.Data, e.DataFormat, &r.Params)
	if err != nil {
		panic(err)
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
	case "json":
		err = json.Unmarshal(data, out)
	default:
		err = json.Unmarshal(data, out)
	}

	return
}
