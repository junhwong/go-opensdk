package alipay

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/junhwong/go-opensdk/common"
)

type Executor struct {
	common.BaseExecutor
	// params2 common.Params
	// params       map[string]string
	// bizContent   map[string]string
	// executed     bool
	// err          error
	client *Client
	// body         []byte
	// resultsField string
}

// func (e *Executor) SetReturnURL(url string) {

// }

// //SetNotifyURL 设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
// func (e *Executor) SetNotifyURL(url string) {

// }

// //SetAppAuthToken 设置第三方应用授权码。接口文档：https://docs.open.alipay.com/20160728150111277227/intro
// func (e *Executor) SetAppAuthToken(url string) {

// }

// func (e *Executor) Results(verifySign bool, outBinding interface{}) (res map[string]interface{}, err error) {
// 	if e.checkExecute() != nil {
// 		err = e.err
// 		return
// 	}
// 	fmt.Println(string(toUTF8Data(e.body)))
// 	err = json.Unmarshal(toUTF8Data(e.body), &res)
// 	if err != nil {
// 		return
// 	}
// 	if v, ok := res[e.resultsField]; ok && v != nil {
// 		res = res[e.resultsField].(map[string]interface{})
// 		if v, ok := res["code"]; ok && v != "10000" {
// 			err = fmt.Errorf("%v", res)
// 			return
// 		}
// 	}
// 	if err == nil && outBinding != nil {
// 		MapToStruct(outBinding, res)
// 	}
// 	return
// }

// //MustGet 网关返回码为 10000 时返回结果，否则错误。错误码参见：https://docs.open.alipay.com/common/105806
// func (e *Executor) MustGet() {
// 	if e.checkExecute() != nil {
// 		panic(e.err)
// 	}
// 	m := make(map[string]interface{})
// 	err := json.Unmarshal(toUTF8Data(e.body), &m)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%+v\n", m)
// 	fmt.Printf("%+v\n", string(toUTF8Data(e.body)))

// }

// //VerifySign 验证相应结果签名。
// //注：一般同步响应支付宝建议可以不用验证。
// func (e *Executor) VerifySign() {

// }

// func (e *Executor) checkExecute() error {
// 	if e.err != nil || e.executed {
// 		return e.err
// 	}
// 	e.body, e.err = e.execute()
// 	return e.err
// }

// func (e *Executor) execute() ([]byte, error) {
// 	params := e.params
// 	s := MapToSortString(params, true, false)

// 	if sign, err := sha256WithRSA(s, e.client.PrivateKey); err != nil {
// 		return nil, err
// 	} else {
// 		params["sign"] = sign
// 	}

// 	urlencode := url.Values{}
// 	for k, v := range params {
// 		urlencode.Add(k, v)
// 	}

// 	res, err := http.Post(e.client.Gateway, "application/x-www-form-urlencoded", strings.NewReader(urlencode.Encode()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	//reader, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
// 	return ioutil.ReadAll(res.Body)
// }

func (e *Executor) Execute(verifySign ...bool) common.Results {
	delete(e.Params, "sign")
	ret := &common.DefaultResults{
		Params: common.Params{},
	}
	params := e.Params.Sort()
	sign, err := sha256WithRSA(params.ToURLParams(false), e.client.PrivateKey)
	if err != nil {
		ret.Err = err
		return ret
	}
	params = append(params, [2]string{"sign", sign})
	u, _ := url.Parse(e.client.Gateway)
	ret.Data, _, ret.Err = e.DoPost(u, strings.NewReader(params.ToURLParams(true)))
	if ret.Err != nil {
		return ret
	}
	ret.ResultString, ret.SignString = extract(ret.Data)
	//签名验证
	if len(verifySign) > 0 && verifySign[0] {
		ret.Err = verifyRSA2(ret.ResultString, ret.SignString, e.client.PublicKey)
		if ret.Err != nil {
			return ret
		}
	}
	reader := toUTF8([]byte(ret.ResultString)) // 支付宝返回编码是GBK，不管传递参数是不是GBK。这是BUG?
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		ret.Err = err
		return ret
	}
	ret.Err = json.Unmarshal(data, &ret.Params)
	return ret
}

// func (e *Executor) DoPost(url string, body io.Reader, tye ...string) ([]byte, *http.Response, error) {
// 	res, err := http.Post(e.client.Gateway, "application/x-www-form-urlencoded", body)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer res.Body.Close()
// 	//reader, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
// 	data, err := ioutil.ReadAll(res.Body)
// 	return data, res, err
// }

// type ResponseResults struct {
// 	common.Params
// 	Err          error
// 	Data         []byte
// 	ResultString string
// 	SignString   string
// }

// func (r *ResponseResults) Set(key string, v interface{}, err error) common.Results {
// 	return r
// }
// func (r *ResponseResults) Get(key string) utils.Converter {
// 	return utils.Convert(r.Params[key], r.Err, true)
// }
// func (r *ResponseResults) Bind(v interface{}, applyType ...string) error {
// 	return nil
// }
// func (r *ResponseResults) Error() error {
// 	return r.Err
// }
// func (r *ResponseResults) Body() []byte {
// 	return r.Data
// }
// func (r *ResponseResults) Success() bool {
// 	return r.Code() == "10000"
// }

// func (r *ResponseResults) Code() string {
// 	return r.Get("code").String()
// }
// func (r *ResponseResults) Message() string {
// 	return r.Get("msg").String()
// }
// func (r *ResponseResults) SubCode() string {
// 	return r.Get("sub_code").String()
// }
// func (r *ResponseResults) SubMessage() string {
// 	return r.Get("sub_msg").String()
// }
// func (r *ResponseResults) Sign() string {
// 	return r.SignString
// }
