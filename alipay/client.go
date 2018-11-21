package alipay

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	"github.com/junhwong/go-opensdk/common"
	"github.com/junhwong/go-utils/crypto"
)

// 待签名样例
// app_id=2014072300007148&biz_content={"button":[{"actionParam":"ZFB_HFCZ","actionType":"out","name":"话费充值"},{"name":"查询","subButton":[{"actionParam":"ZFB_YECX","actionType":"out","name":"余额查询"},{"actionParam":"ZFB_LLCX","actionType":"out","name":"流量查询"},{"actionParam":"ZFB_HFCX","actionType":"out","name":"话费查询"}]},{"actionParam":"http://m.alipay.com","actionType":"link","name":"最新优惠"}]}&charset=GBK&method=alipay.mobile.public.menu.add&sign_type=RSA2&timestamp=2014-07-24 03:07:50&version=1.0

// type CommonResponse struct {
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

// func BuildParams(method string) map[string]string {
// 	params := map[string]string{}
// 	params["method"] = method
// 	return params
// }

// type OauthTokenResponse struct {
// }

func (c *Client) Build(methodOrURL string, params ...common.Params) common.Executor {
	if len(params) > 1 {
		panic("无效的参数")
	}
	var p common.Params
	if len(params) == 1 {
		p = params[0]
	} else {
		p = common.Params{}
	}
	p["method"] = methodOrURL
	p["app_id"] = c.AppID
	p["format"] = "json"
	p["charset"] = "utf-8"
	p["version"] = "1.0"
	p["sign_type"] = "RSA2"
	p["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	return &Executor{
		BaseExecutor: common.BaseExecutor{
			Params: p,
		},
		client: c,
	}
}
