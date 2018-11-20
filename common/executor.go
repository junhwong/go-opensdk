package common

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/junhwong/go-utils"
)

var (
	ErrNotExecuteYet = errors.New("not execute yet.")
)

type ResponseParseFn func(*http.Response) *Response

// Executor 用于执行请求的相关上下文。
// 每个 Executor 都可以被多次执行(注意：接口业务是否允许)。
type Executor interface {
	// Execute 执行请求并返回结果。
	// verifySign 表示是否需要同步验证签名，默认 `false` 不验证。
	Execute(verifySign ...bool) *Response
}

type DefaultExecutor struct {
	Params      Parameters
	executed    bool
	executedErr error
	//client       *Client
	body             []byte
	ParseFn          ResponseParseFn
	VerifySignFn     func(*Response) bool
	RequestURL       *url.URL
	RequestMethod    string
	RequestEncoding  string
	ResponseEncoding string
}

// Results 接口业务执行成功的结果。
type Results interface {
	// Set 重新设置结果中的某个值，以用于数据绑定。
	Set(key string, v interface{}, err error) Results
	// Get 获取结果中的某个值。
	Get(key string) utils.Converter
	// Bind 绑定到一个结构体。
	// applyType 绑定是字段匹配的类型, 值：default、json、xml。
	Bind(v interface{}, applyType ...string) error
}

// Response 执行接口返回的响应。
type Response struct {
	Err        error
	Sign       string `json:"sign"`
	Code       string `json:"code"`
	Message    string `json:"msg"`
	SubCode    string `json:"sub_code"`
	SubMessage string `json:"sub_msg"`
	Results    Results
}

func (e *DefaultExecutor) SetReturnURL(url string) {

}

//SetNotifyURL 设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
func (e *DefaultExecutor) SetNotifyURL(url string) {

}

//SetAppAuthToken 设置第三方应用授权码。接口文档：https://docs.open.alipay.com/20160728150111277227/intro
func (e *DefaultExecutor) SetAppAuthToken(url string) {

}

/*
{"alipay_system_oauth_token_response":{"access_token":"authusrBe456e5dcbf714cebacd620d83c047F79","alipay_user_id":"20881041302782186927417731510379","expires_in":1296000,"re_expires_in":2592000,"refresh_token":"authusrB361200341ef14df8b7eee42682b31X79","user_id":"2088002463165793"},"sign":"F/A0Hmg88ven3owWJ0umE9n9iMQYB+QoHLF/hd+BMcNoSzA4gNqr07BDWSnHVQPwQiPf3JQuj+b+4RTpPAh07FVPvbjdkIBsN84dwwejyZ8pwQlV/CqZEuNrrbFwGadnFDjtKmppj5qDh6YcHDei6TEvtxQY2Uz1ZxlxvBdqFRO+lyEiefgaw4ZaD1B3ccWrueu6pQiqgm6h/23//N6hEgCKzI3rJdvCLKVjehTaolYWdbVBnnDaxMEISEmGDOPbjMBfaY7YT/eiJ4I9XktRTtIvvJdMVpCxC6mgMnFP3szFlbSC7bj8o6l+z6CDtqDPCzfklJ2LKlYxom7vV8Y+/Q=="}

{"alipay_system_oauth_token_response":{"access_token":"authusrBe456e5dcbf714cebacd620d83c047F79","alipay_user_id":"20881041302782186927417731510379","expires_in":1296000,"re_expires_in":2592000,"refresh_token":"authusrB361200341ef14df8b7eee42682b31X79","user_id":"2088002463165793"},"sign":"F/A0Hmg88ven3owWJ0umE9n9iMQYB+QoHLF/hd+BMcNoSzA4gNqr07BDWSnHVQPwQiPf3JQuj+b+4RTpPAh07FVPvbjdkIBsN84dwwejyZ8pwQlV/CqZEuNrrbFwGadnFDjtKmppj5qDh6YcHDei6TEvtxQY2Uz1ZxlxvBdqFRO+lyEiefgaw4ZaD1B3ccWrueu6pQiqgm6h/23//N6hEgCKzI3rJdvCLKVjehTaolYWdbVBnnDaxMEISEmGDOPbjMBfaY7YT/eiJ4I9XktRTtIvvJdMVpCxC6mgMnFP3szFlbSC7bj8o6l+z6CDtqDPCzfklJ2LKlYxom7vV8Y+/Q=="}
{"alipay_user_info_share_response":{"code":"20001","msg":"Insufficient Token Permissions","sub_code":"aop.invalid-auth-token","sub_msg":"无效的访问令牌"},"sign":"fFtejz7s352L9iWGNp1/aSqG33oXGkt2Syj+Ik3zx0Qh2XkR+Lb/9O0OJ87VV4/WUQA21g8fVgyYSWPbv4dksIpJV55ubRO01LnyLddBBYzXuIph+WjTcWWk0OZktxqqAKvTu1zn65REwLSVNxVISL/KqhpesLXesMX3y84dlq3vgQ0AiVp7aG8+q7xGP6Jb4NwHZA6eY9RERNCRbNxYKZ+57CM85J6HxpebndzSCiyrnkCN+teOoiuk1ICrmEdDqLmF+A25SvroIOfrmMkGZTp1RLoVDkYdE/3M4B2YSxaXJqexpGsqkYStEfMkTV0CDA9GDs3hbfv392LSzjoTbQ=="}

*/

// Execute 执行请求并返回结果。
// verifySign 表示是否需要同步验证签名，默认 `false` 不验证。
func (e *DefaultExecutor) Execute(verifySign ...bool) *Response {
	switch strings.ToUpper(e.RequestMethod) {
	case "POST":
		return e.DoPost()
	case "GET":
		return e.DoGet()
	default:
		panic(fmt.Errorf("暂不支持请求方法: %s", e.RequestMethod))
	}
}

func (e *DefaultExecutor) DoPost() *Response {
	return e.ParseFn(nil)
}
func (e *DefaultExecutor) DoGet() *Response {
	return e.ParseFn(nil)
}

//VerifySign 验证相应结果签名。
//注：一般同步响应支付宝建议可以不用验证。
func (e *DefaultExecutor) VerifySign() {

}

func (e *DefaultExecutor) checkExecute() error {
	// if e.err != nil || e.executed {
	// 	return e.err
	// }
	// e.body, e.err = e.execute()
	// return e.err
	return nil
}

func (e *DefaultExecutor) execute() ([]byte, error) {
	// params := e.params
	// s := MapToSortString(params, true, false)

	// if sign, err := sha256WithRSA(s, e.client.PrivateKey); err != nil {
	// 	return nil, err
	// } else {
	// 	params["sign"] = sign
	// }

	// urlencode := url.Values{}
	// for k, v := range params {
	// 	urlencode.Add(k, v)
	// }

	// res, err := http.Post(e.client.Gateway, "application/x-www-form-urlencoded", strings.NewReader(urlencode.Encode()))
	// if err != nil {
	// 	return nil, err
	// }
	// defer res.Body.Close()
	// //reader, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
	// return ioutil.ReadAll(res.Body)
	return nil, nil
}
