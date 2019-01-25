package wechat

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/junhwong/go-opensdk/opensdk"
)

type WechatClient struct {
	AppID  string //微信分配的小程序ID
	Secret string
}

type WechatPayClient struct {
	WechatClient
	APICertFile string
	MchID       string // 微信支付分配的商户号
	MchKey      string
	tlsClient   *http.Client
}

func NewClient(appid, secret string) *WechatClient {
	return &WechatClient{
		AppID:  appid,
		Secret: secret,
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
	e.Request = c.doRequest
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
func (c *WechatPayClient) VerifySign(params, signType string) (bool, error) {
	return false, nil
}

func (c *WechatPayClient) getHttpClient(useTLS bool) (*http.Client, error) {
	if !useTLS {
		return http.DefaultClient, nil
	}
	if c.tlsClient == nil {
		certData, err := ioutil.ReadFile(c.APICertFile)
		if err != nil {
			return nil, err
		}
		// 将pkcs12证书转成pem
		cert := pkcs12ToPem(certData, c.MchID)
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
			DisableCompression: true,
		}
		c.tlsClient = &http.Client{Transport: transport}
	}
	return c.tlsClient, nil
}

func (c *WechatPayClient) doRequest(def *opensdk.DefaultExecutor) (response *http.Response, requestLog string, err error) {
	log := ""
	hc, err := c.getHttpClient(def.TLS)
	if err != nil {
		return nil, log, err
	}
	signType := def.Get("sign_type").String()
	// delete(def.Params, "sign_type")
	params := def.Params.Sort()
	log += "URL:" + def.APIURL
	log += "\nsign params:" + params.ToURLParams()
	sign, err := c.Sign(params.ToURLParams(), signType)
	if err != nil {
		return nil, log, err
	}
	// delete(def.Params, "sign_type")
	// params = def.Params.Sort()
	params.Append("sign", sign)
	body := params.ToXML()
	log += "\n"
	log += "body:" + body
	log += "\n"
	resp, err := hc.Post(def.APIURL, "text/xml", strings.NewReader(body)) // TODO: HTTP METHOD
	return resp, log, err
}
