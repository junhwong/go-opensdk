package alipay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Executor struct {
	params       map[string]string
	bizContent   map[string]string
	executed     bool
	err          error
	client       *Client
	body         []byte
	resultsField string
}

func (e *Executor) SetReturnURL(url string) {

}

//SetNotifyURL 设置支付宝服务器主动通知商户服务器里指定的页面http/https路径。
func (e *Executor) SetNotifyURL(url string) {

}

//SetAppAuthToken 设置第三方应用授权码。接口文档：https://docs.open.alipay.com/20160728150111277227/intro
func (e *Executor) SetAppAuthToken(url string) {

}

func (e *Executor) Results(verifySign bool, outBinding interface{}) (res map[string]interface{}, err error) {
	if e.checkExecute() != nil {
		err = e.err
		return
	}
	fmt.Println(string(toUTF8Data(e.body)))
	err = json.Unmarshal(toUTF8Data(e.body), &res)
	if err != nil {
		return
	}
	if v, ok := res[e.resultsField]; ok && v != nil {
		res = res[e.resultsField].(map[string]interface{})
		if v, ok := res["code"]; ok && v != "10000" {
			err = fmt.Errorf("%v", res)
			return
		}
	}
	if err == nil && outBinding != nil {
		MapToStruct(outBinding, res)
	}
	return
}

//MustGet 网关返回码为 10000 时返回结果，否则错误。错误码参见：https://docs.open.alipay.com/common/105806
func (e *Executor) MustGet() {
	if e.checkExecute() != nil {
		panic(e.err)
	}
	m := make(map[string]interface{})
	err := json.Unmarshal(toUTF8Data(e.body), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", m)
	fmt.Printf("%+v\n", m["error_response"])

}

//VerifySign 验证相应结果签名。
//注：一般同步响应支付宝建议可以不用验证。
func (e *Executor) VerifySign() {

}

func (e *Executor) checkExecute() error {
	if e.err != nil || e.executed {
		return e.err
	}
	e.body, e.err = e.execute()
	return e.err
}

func (e *Executor) execute() ([]byte, error) {
	params := e.params
	s := MapToSortString(params, true, false)

	if sign, err := sha256WithRSA(s, e.client.PrivateKey); err != nil {
		return nil, err
	} else {
		params["sign"] = sign
	}

	urlencode := url.Values{}
	for k, v := range params {
		urlencode.Add(k, v)
	}

	res, err := http.Post(e.client.Gateway, "application/x-www-form-urlencoded", strings.NewReader(urlencode.Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	//reader, _ := iconv.NewReader(res.Body, "gb2312", "utf-8")
	return ioutil.ReadAll(res.Body)
}
