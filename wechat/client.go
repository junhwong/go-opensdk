package wechat

import (
	"net/http"
	"strings"

	logs "github.com/junhwong/go-logs"
	"github.com/junhwong/go-opensdk/opensdk"
)

type WechatClient struct {
	opensdk.ClientBase
	AppID  string //微信分配的小程序ID
	Secret string
}

func NewClient(appid, secret string) *WechatClient {
	return &WechatClient{
		ClientBase: opensdk.ClientBase{},
		AppID:      appid,
		Secret:     secret,
	}
}

type WechatPayClient struct {
	WechatClient
	APICertFile string // 过时
	ServiceID   string
	MchID       string // 微信支付分配的商户号
	MchKey      string
}

func NewPayClient(appid, secret, mchID, mchKey, serviceID string) *WechatPayClient {
	return &WechatPayClient{
		WechatClient: *NewClient(appid, secret),
		ServiceID:    serviceID,
		MchID:        mchID,
		MchKey:       mchKey,
	}
}

func (c *WechatPayClient) buildPOST(url string, params opensdk.Params) opensdk.Executor {
	params["appid"] = c.AppID
	params["mch_id"] = c.MchID
	params["nonce_str"] = opensdk.RandomString(10)
	params["sign_type"] = "MD5"

	return c.BuildExecutor(url, params)
}

func (c *WechatPayClient) BuildExecutor(url string, params opensdk.Params) opensdk.Executor {
	e := opensdk.DefaultExecutor{
		Params:     params,
		Client:     c,
		HTTPMethod: "POST",
		APIURL:     url,
	}
	e.BuildRequest = c.doRequest
	return e.ResultValidator(func(p opensdk.Params) (ok bool, code string, msg string, subcode string, submsg string) {
		code = p.Get("return_code").String()
		msg = p.Get("return_msg").String()
		subcode = p.Get("err_code").String()
		submsg = p.Get("err_code_des").String()

		if code == "SUCCESS" && p.Get("result_code").String() == "SUCCESS" {
			ok = true
		}
		return
	})
}

func (c *WechatPayClient) Sign(params, signType string) (string, error) {
	switch signType {
	case "HMAC-SHA256":
		return opensdk.Sha256Hmac(params, nil) // TODO: key
	case "SHA1":
		return opensdk.Sha1(params)
	default:
		return opensdk.MD5(params + "&key=" + c.MchKey)
	}
	// return "", errors.New("签名算法不支持：" + signType)
}

func (c *WechatClient) VerifySign(params, signType string) (bool, error) {
	return false, nil
}

func (c *WechatPayClient) doRequest(def *opensdk.DefaultExecutor) (req *http.Request, err error) {
	log := ""
	signType := def.Get("sign_type").String()
	// delete(def.Params, "sign_type")
	params := def.Params.Sort()
	log += "URL:" + def.APIURL
	log += "\nsign params:" + params.ToURLParams()
	sign, err := c.Sign(params.ToURLParams(), signType)
	if err != nil {
		return nil, err
	}
	// delete(def.Params, "sign_type")
	// params = def.Params.Sort()
	params.Append("sign", sign)
	body := params.ToXML()
	log += "\n"
	log += "body:" + body
	log += "\n"
	logs.Prefix("go-opensdk.wechat").Debug("request params:", log, params.ToURLParams(false))
	req, err = http.NewRequest("POST", def.APIURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml")
	return req, nil
}
