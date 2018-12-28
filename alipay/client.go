package alipay

import (
	"crypto/rsa"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/junhwong/go-opensdk/opensdk"
	"github.com/junhwong/go-utils/crypto"
)

// 待签名样例
// app_id=2014072300007148&biz_content={"button":[{"actionParam":"ZFB_HFCZ","actionType":"out","name":"话费充值"},{"name":"查询","subButton":[{"actionParam":"ZFB_YECX","actionType":"out","name":"余额查询"},{"actionParam":"ZFB_LLCX","actionType":"out","name":"流量查询"},{"actionParam":"ZFB_HFCX","actionType":"out","name":"话费查询"}]},{"actionParam":"http://m.alipay.com","actionType":"link","name":"最新优惠"}]}&charset=GBK&method=alipay.mobile.public.menu.add&sign_type=RSA2&timestamp=2014-07-24 03:07:50&version=1.0

// type opensdkResponse struct {
// 	Sign       string `json:"sign"`
// 	Code       string `json:"code"`
// 	Message    string `json:"msg"`
// 	SubCode    string `json:"sub_code"`
// 	SubMessage string `json:"sub_msg"`
// }

type Client struct {
	AppID      string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Gateway    string
}

func NewClient(gateway, appID, privateKeyPath, alipayPublicKeyPath string) *Client {
	c := &Client{
		AppID:   appID,
		Gateway: gateway,
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

func (c *Client) Build(methodOrURL string, params opensdk.Params) opensdk.Executor {

	params["method"] = methodOrURL
	params["app_id"] = c.AppID
	params["format"] = "json"
	params["charset"] = "utf-8"
	params["version"] = "1.0"
	params["sign_type"] = "RSA2"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	e := &opensdk.DefaultExecutor{
		Params:     params,
		Client:     c,
		APIURL:     c.Gateway,
		HTTPMethod: "POST",
	}
	e.Request = c.doRequest
	return e.ResultValidator(func(p opensdk.Params) (ok bool, code string, msg string, subcode string, submsg string) {
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
	return opensdk.Sha256RSA(params, c.PrivateKey)
}

// VerifySign 验签
func (c *Client) VerifySign(params, signType string) (bool, error) {
	return false, nil
}

func (c *Client) doRequest(def *opensdk.DefaultExecutor) (response *http.Response, requestLog string, err error) {
	signType := def.Get("sign_type").String()
	params := def.Params.Sort()
	log := "sign params:" + params.ToURLParams()
	sign, err := c.Sign(params.ToURLParams(), signType)
	if err != nil {
		return nil, log, err
	}
	params.Append("sign", sign)
	body := params.ToURLParams(true)
	log += "\nbody:" + body
	def.Decoder = func(data []byte, dataFormat string, out *opensdk.Params) (err error) {
		dataStr, signStr := extract(data) // 分离结果和签名
		(*out)["sign"] = signStr

		reader := toUTF8([]byte(dataStr)) // 支付宝返回编码是GBK，不管传递参数是不是GBK。是BUG?
		newData, err := ioutil.ReadAll(reader)
		if err != nil {
			if err != io.EOF {
				return err
			}
		}
		return opensdk.DefaultDecoder(newData, dataFormat, out)
	}
	resp, err := http.Post(def.APIURL, "application/x-www-form-urlencoded", strings.NewReader(body)) // TODO: HTTP METHOD
	return resp, log, err
}
