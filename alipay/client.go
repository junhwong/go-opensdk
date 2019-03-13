package alipay

import (
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/junhwong/go-logs"
	"github.com/junhwong/go-opensdk/core"
	"github.com/junhwong/go-utils/crypto"
)

// 待签名样例
// app_id=2014072300007148&biz_content={"button":[{"actionParam":"ZFB_HFCZ","actionType":"out","name":"话费充值"},{"name":"查询","subButton":[{"actionParam":"ZFB_YECX","actionType":"out","name":"余额查询"},{"actionParam":"ZFB_LLCX","actionType":"out","name":"流量查询"},{"actionParam":"ZFB_HFCX","actionType":"out","name":"话费查询"}]},{"actionParam":"http://m.alipay.com","actionType":"link","name":"最新优惠"}]}&charset=GBK&method=alipay.mobile.public.menu.add&sign_type=RSA2&timestamp=2014-07-24 03:07:50&version=1.0

// type coreResponse struct {
// 	Sign       string `json:"sign"`
// 	Code       string `json:"code"`
// 	Message    string `json:"msg"`
// 	SubCode    string `json:"sub_code"`
// 	SubMessage string `json:"sub_msg"`
// }

type Client struct {
	core.ClientBase
	AppID          string
	PrivateKey     *rsa.PrivateKey
	PublicKey      *rsa.PublicKey
	Gateway        string
	httpClientFunc func(useTLS bool) (*http.Client, error)
}

func (c *Client) SetHttpClientFunc(fn func(twowayAuthentication bool) (*http.Client, error)) {
	c.httpClientFunc = fn
}
func NewClient(gateway, appID, privateKeyPath, alipayPublicKeyPath string) *Client {
	c := &Client{
		ClientBase: core.ClientBase{},
		AppID:      appID,
		Gateway:    gateway,
	}

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	if key, err := crypto.PKCS1PrivateKey(privateKey); err != nil {
		panic(err)
	} else {
		c.PrivateKey = key
	}

	publicKey, err := ioutil.ReadFile(alipayPublicKeyPath)
	if err != nil {
		panic(err)
	}
	if key, err := crypto.PKCS1PublicKey(publicKey); err != nil {
		panic(err)
	} else {
		c.PublicKey = key
	}

	return c
}

// 填充通用参数
func (c *Client) fillParams(method string, params core.Params) {
	params["sdk"] = "go-opensdk" // TODO: 版本号
	params["method"] = method
	params["app_id"] = c.AppID
	params["format"] = "json"
	params["charset"] = "utf-8"
	params["version"] = "1.0"
	params["sign_type"] = "RSA2"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
}

func (c *Client) Build(methodOrURL string, params core.Params) core.Executor {

	c.fillParams(methodOrURL, params)

	e := &core.DefaultExecutor{
		Params:     params,
		Client:     c,
		APIURL:     c.Gateway,
		HTTPMethod: "POST",
	}
	e.BuildRequest = c.doRequest
	return e.ResultValidator(func(p core.Params) (ok bool, code string, msg string, subcode string, submsg string) {
		code = p.Get("code").String()
		msg = p.Get("msg").String()
		subcode = p.Get("sub_code").String()
		submsg = p.Get("sub_msg").String()

		if code == "10000" {
			ok = true
		}
		return
	})
}

func (c *Client) buildWithBizContent(method string, biz core.Params) core.Executor {
	params := core.Params{}
	c.fillParams(method, params)
	params["biz_content"] = biz
	e := &core.DefaultExecutor{
		Params:     params,
		Client:     c,
		APIURL:     c.Gateway,
		HTTPMethod: "POST",
	}
	e.BuildRequest = c.doRequest
	return e.ResultValidator(func(p core.Params) (ok bool, code string, msg string, subcode string, submsg string) {
		code = p.Get("code").String()
		msg = p.Get("msg").String()
		subcode = p.Get("sub_code").String()
		submsg = p.Get("sub_msg").String()

		if code == "10000" {
			ok = true
		}
		return
	})
}

// Sign 签名
func (c *Client) Sign(params, signType string) (string, error) {
	return core.Sha256RSA(params, c.PrivateKey)
}

// VerifySign 验签
func (c *Client) VerifySign(params, signType string) (bool, error) {
	return false, nil
}

func (c *Client) doRequest(def *core.DefaultExecutor) (response *http.Request, err error) {
	log := ""
	signType := def.Get("sign_type").String()
	params := def.Params.Sort()
	log = "sign params:" + params.ToURLParams()
	sign, err := c.Sign(params.ToURLParams(), signType)
	if err != nil {
		return nil, err
	}
	params.Append("sign", sign)
	body := params.ToURLParams(true)
	log += "\nbody:" + body
	logs.Prefix("go-core.alipay").Debug("request params:", log, params.ToURLParams(false))
	def.Decoder = func(data []byte, dataFormat string, out *core.Params) (err error) {
		dataStr, signStr := extract(data) // 分离结果和签名
		(*out)["sign"] = signStr

		newData, err := core.ToUTF8Data([]byte(dataStr)) // 支付宝返回编码是GBK，不管传递参数是不是GBK。是BUG?
		// newData, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		return core.DefaultDecoder(newData, dataFormat, out)
	}
	req, err := http.NewRequest("POST", def.APIURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}
