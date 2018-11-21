package common

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

	SetReturnURL(url string)

	//SetNotifyURL 设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	SetNotifyURL(url string)
}

type BaseExecutor struct {
	Params Params
	// params       map[string]string
	// bizContent   map[string]string
	// executed     bool
	// err          error
	// client *Client
	// body         []byte
	// resultsField string
}

func (e *BaseExecutor) SetReturnURL(url string) {

}

//SetNotifyURL 设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
func (e *BaseExecutor) SetNotifyURL(url string) {

}

//SetAppAuthToken 设置第三方应用授权码。接口文档：https://docs.open.alipay.com/20160728150111277227/intro
func (e *BaseExecutor) SetAppAuthToken(url string) {

}

func (e *BaseExecutor) DoPost(url *url.URL, body io.Reader, tye ...string) ([]byte, *http.Response, error) {
	res, err := http.Post(url.String(), "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	//reader, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
	data, err := ioutil.ReadAll(res.Body)
	return data, res, err
}

// type DefaultExecutor struct {
// 	Params      Params
// 	executed    bool
// 	executedErr error
// 	//client       *Client
// 	body []byte
// 	// ParseFn          ResponseParseFn
// 	VerifySignFn     func(Results) bool
// 	RequestURL       *url.URL
// 	RequestMethod    string
// 	RequestEncoding  string
// 	ResponseEncoding string
// }

// func (e *DefaultExecutor) SetReturnURL(url string) {

// }

// //SetNotifyURL 设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
// func (e *DefaultExecutor) SetNotifyURL(url string) {

// }

// //SetAppAuthToken 设置第三方应用授权码。接口文档：https://docs.open.alipay.com/20160728150111277227/intro
// func (e *DefaultExecutor) SetAppAuthToken(url string) {

// }

// // Execute 执行请求并返回结果。
// // verifySign 表示是否需要同步验证签名，默认 `false` 不验证。
// func (e *DefaultExecutor) Execute(verifySign ...bool) Results {
// 	switch strings.ToUpper(e.RequestMethod) {
// 	case "POST":
// 		return e.DoPost()
// 	case "GET":
// 		return e.DoGet()
// 	default:
// 		panic(fmt.Errorf("暂不支持请求方法: %s", e.RequestMethod))
// 	}
// }

// func (e *DefaultExecutor) DoPost() Results {
// 	return nil
// }
// func (e *DefaultExecutor) DoGet() Results {
// 	return nil
// }

// //VerifySign 验证相应结果签名。
// //注：一般同步响应支付宝建议可以不用验证。
// func (e *DefaultExecutor) VerifySign() {

// }

// func (e *DefaultExecutor) checkExecute() error {
// 	// if e.err != nil || e.executed {
// 	// 	return e.err
// 	// }
// 	// e.body, e.err = e.execute()
// 	// return e.err
// 	return nil
// }

// func (e *DefaultExecutor) execute() ([]byte, error) {
// 	// params := e.params
// 	// s := MapToSortString(params, true, false)

// 	// if sign, err := sha256WithRSA(s, e.client.PrivateKey); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	params["sign"] = sign
// 	// }

// 	// urlencode := url.Values{}
// 	// for k, v := range params {
// 	// 	urlencode.Add(k, v)
// 	// }

// 	// res, err := http.Post(e.client.Gateway, "application/x-www-form-urlencoded", strings.NewReader(urlencode.Encode()))
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// defer res.Body.Close()
// 	// //reader, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
// 	// return ioutil.ReadAll(res.Body)
// 	return nil, nil
// }
