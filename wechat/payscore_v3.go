package wechat

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	logs "github.com/junhwong/go-logs"
	"github.com/junhwong/go-opensdk/opensdk"
)

func (c *WechatPayClient) PayScoreV3() *WechatPayScoreV3 {
	return &WechatPayScoreV3{client: c}
}

type WechatPayScoreV3 struct {
	client *WechatPayClient
}

func (c *WechatPayScoreV3) Parent() *WechatPayClient {
	return c.client
}
func (c *WechatPayScoreV3) buildExecutor(url, method string, params opensdk.Params) opensdk.Executor {
	e := opensdk.DefaultExecutor{
		Params:     params,
		Client:     c.client,
		HTTPMethod: method,
		APIURL:     url,
	}
	e.BuildRequest = c.buildRequest

	return e.ResultValidator(c.getResultValidator)
}

func (c *WechatPayScoreV3) buildRequest(e *opensdk.DefaultExecutor) (req *http.Request, err error) {
	body := ""
	rawurl := e.APIURL
	if e.HTTPMethod == "GET" {
		if len(e.Params) > 0 {
			rawurl += "?"
			rawurl += e.Params.Sort().ToURLParams(true)
		}
	} else {
		body = e.Params.Sort2(true).ToJSON2(true)
	}

	req, err = http.NewRequest(e.HTTPMethod, rawurl, strings.NewReader(body))
	if err != nil {
		return
	}
	//https://github.com/wechatpay-apiv3/wechatpay-apache-httpclient/blob/master/src/main/java/com/wechat/pay/contrib/apache/httpclient/auth/PrivateKeySigner.java
	//https://github.com/wechatpay-apiv3/wechatpay-apache-httpclient/blob/master/src/main/java/com/wechat/pay/contrib/apache/httpclient/auth/WechatPay2Credentials.java

	u, _ := url.Parse(rawurl)
	reqPath := u.EscapedPath()
	if u.RawQuery != "" {
		reqPath += "?"
		reqPath += u.RawQuery
	}
	authp := opensdk.Params{
		"nonce_str": opensdk.RandomString(10),
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"serial_no": c.client.APICertSerialNo,
		"mchid":     c.client.MchID,
		"signature": "signature",
	}
	msg := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		e.HTTPMethod,       // method
		reqPath,            // url
		authp["timestamp"], // timestamp
		authp["nonce_str"], // nonce_str
		body,               // body
	)
	logs.Prefix("go-opensdk.wechat").Debug("signature message: ", msg)
	signature, err := opensdk.Sha256RSA(msg, c.client.APICert.PrivateKey.(*rsa.PrivateKey))
	if err != nil {
		return nil, err
	}
	authp["signature"] = signature
	authv := "WECHATPAY2-SHA256-RSA2048 " + authp.Sort().ToHeaderParams()
	req.Header.Add("Authorization", authv)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	// logs.Prefix("go-opensdk.wechat").Debug("Authorization: ", authv)
	return req, nil
}
func (c *WechatPayScoreV3) getResultValidator(p opensdk.Params) (ok bool, code string, msg string, subcode string, submsg string) {
	status := p.Get("response.StatusCode").Int()

	code = p.Get("code").String()
	msg = p.Get("message").String()
	subcode = code // p.Get("err_code").String()
	submsg = msg   // p.Get("err_code_des").String()

	ok = status == 200 || status == 204

	if ok && code == "" {
		code = "OK"
	}
	return
}

// ServiceOrder 用于查询单笔微信支付分订单详细信息。
// [接口文档](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/payscore/chapter3_2.shtml)
func (p *WechatPayScoreV3) ServiceOrderQuery(ctx context.Context, out_order_no string, options ...ParamsOption) opensdk.Results {
	// params.Get("out_order_no")
	//https://wechatpay-api.gitbook.io/wechatpay-api-v3/qian-ming-zhi-nan-1/qian-ming-sheng-cheng
	params := opensdk.Params{
		"out_order_no": out_order_no,
		"appid":        p.client.AppID,
		"service_id":   p.client.ServiceID,
	}
	for _, it := range options {
		if it == nil {
			continue
		}
		it(params)
	}

	executor := p.buildExecutor("https://api.mch.weixin.qq.com/v3/payscore/serviceorder",
		"GET", params)
	return executor.Execute(false)
}

// ServiceOrderCreate 申请创建微信支付分订单。
// [接口文档](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/payscore/chapter3_1.shtml)
func (p *WechatPayScoreV3) ServiceOrderCreate(ctx context.Context, biz PayScoreRentParams, options ...ParamsOption) opensdk.Results {
	params := opensdk.Params{
		"out_order_no":         biz.OutOrderNo,
		"appid":                p.client.AppID,
		"service_id":           p.client.ServiceID,
		"service_introduction": biz.GoodsName,
		"time_range": opensdk.Params{
			"start_time": "OnAccept", //biz.StartTime.Format("20060102150405"), // 传入固定值OnAccept表示用户确认订单成功时间为服务开始时间
		},
		"risk_fund": opensdk.Params{
			"name":   "DEPOSIT",
			"amount": biz.DepositAmount,
		},
		"need_user_confirm": true,
		"notify_url":        biz.NotifyURL,
	}
	for _, it := range options {
		if it == nil {
			continue
		}
		it(params)
	}

	executor := p.buildExecutor("https://api.mch.weixin.qq.com/v3/payscore/serviceorder",
		"POST", params)
	//PayScoreRentCreateResults
	return executor.Execute(false)
}

// ServiceOrderCancel 取消支付分订单。
// [接口文档](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/payscore/chapter3_1.shtml)
func (p *WechatPayScoreV3) ServiceOrderCancel(ctx context.Context, biz struct {
	OutOrderNo string //商户服务订单号
	Reason     string //取消原因
}, options ...ParamsOption) opensdk.Results {
	params := opensdk.Params{
		"appid":      p.client.AppID,
		"service_id": p.client.ServiceID,
		"reason":     biz.Reason,
	}
	for _, it := range options {
		if it == nil {
			continue
		}
		it(params)
	}

	executor := p.buildExecutor(fmt.Sprintf("https://api.mch.weixin.qq.com/v3/payscore/serviceorder/%s/cancel", biz.OutOrderNo),
		"POST", params)
	return executor.Execute(false)
}

// ServiceOrderComplete 完结微信支付分订单。用户使用服务完成后，商户可通过此接口完结订单。
// [接口文档](https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/payscore/chapter3_5.shtml)
func (p *WechatPayScoreV3) ServiceOrderComplete(ctx context.Context, biz struct {
	OutOrderNo  string //商户服务订单号
	Title       string //取消原因
	TotalAmount int64
	StartTime   time.Time
	EndTime     time.Time
}, options ...ParamsOption) opensdk.Results {
	params := opensdk.Params{
		"appid":      p.client.AppID,
		"service_id": p.client.ServiceID,
		"time_range": opensdk.Params{
			// "start_time": biz.StartTime.Format("20060102150405"),
			"end_time": biz.EndTime.Format("20060102150405"), // 传入固定值OnAccept表示用户确认订单成功时间为服务开始时间
		},
		"total_amount": biz.TotalAmount,
		"post_payments": []opensdk.Params{
			{
				"name":   biz.Title,
				"amount": biz.TotalAmount,
			},
		},
	}
	for _, it := range options {
		if it == nil {
			continue
		}
		it(params)
	}

	executor := p.buildExecutor(fmt.Sprintf("https://api.mch.weixin.qq.com/v3/payscore/serviceorder/%s/complete", biz.OutOrderNo),
		"POST", params)
	return executor.Execute(false)
}
