package wechat

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/junhwong/go-opensdk/opensdk"
)

func genString() string {
	rand.Seed(time.Now().Unix())
	rnd := strconv.FormatInt(rand.Int63(), 15)
	if len(rnd) > 15 {
		rnd = rnd[:15]
	}
	return fmt.Sprintf("%s", rnd)
}

// UnifiedOrder 在微信支付服务后台生成预支付交易单。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_1
//
// 结果通知：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_7&index=8
func (c *WechatPayClient) UnifiedOrder(openID, outTradeNo, body, spbillIP string, totalFee int64) opensdk.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/unifiedorder", opensdk.Params{
		"body":             body,       // 商品简单描述，该字段请按照规范传递
		"out_trade_no":     outTradeNo, // 商户系统内部订单号
		"fee_type":         "CNY",      //
		"total_fee":        totalFee,   //
		"spbill_create_ip": spbillIP,   // APP和H5支付提交用户端ip，Native支付填调用微信支付API的机器IP。
		"trade_type":       "JSAPI",    //
		"openid":           openID,     // trade_type=JSAPI，此参数必传
	}).UseXML(true)
}

// OrderQuery 微信支付订单的查询。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_2
func (c *WechatPayClient) OrderQuery(transactionID, outTradeNo string) opensdk.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/orderquery", opensdk.Params{
		"transaction_id": transactionID, // 微信的订单号，优先使用
		"out_trade_no":   outTradeNo,    // 商户系统内部订单号
	}).UseXML(true)
}

// CloseOrder 关单接口。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_3
func (c *WechatPayClient) CloseOrder(outTradeNo string) opensdk.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/closeorder", opensdk.Params{
		"out_trade_no": outTradeNo, // 商户系统内部订单号
	}).UseXML(true)
}

// Refund 申请退款。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_3
//
// 结果通知：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_16&index=10
func (c *WechatPayClient) Refund(transactionID, outTradeNo, outRefundNo string, totalFee, refundFee int64) opensdk.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/secapi/pay/refund", opensdk.Params{
		"transaction_id":  transactionID, // 微信的订单号，优先使用
		"out_trade_no":    outTradeNo,    // 商户系统内部订单号
		"out_refund_no":   outRefundNo,   // 商户系统内部的退款单号
		"total_fee":       totalFee,      // 订单总金额
		"refund_fee":      refundFee,     // 退款总金额
		"refund_fee_type": "CNY",         // 货币类型
		"refund_desc":     "refund_desc", // 退款原因
	}).UseXML(true).UseTLS(true)
}

// RefundQuery 查询退款。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_3
func (c *WechatPayClient) RefundQuery(transactionID, outTradeNo, outRefundNo, refundID string) opensdk.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/refundquery", opensdk.Params{
		"transaction_id": transactionID, // 微信的订单号，优先使用
		"out_trade_no":   outTradeNo,    // 商户系统内部订单号
		"out_refund_no":  outRefundNo,   // 商户系统内部的退款单号
		"refund_id":      refundID,      // 退款单号
	}).UseXML(true)
}

//BuildMiniProgramRequestPaymentParams 构建小程序调起支付API参数。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=7_7&index=5
func (c *WechatPayClient) BuildMiniProgramRequestPaymentParams(prepayID string) (params opensdk.Params, err error) {
	params = opensdk.Params{
		"appId":     c.AppID,
		"timeStamp": time.Now().Unix(),
		"nonceStr":  opensdk.RandomString(10),
		"package":   "prepay_id=" + prepayID,
		"signType":  "MD5",
	}
	params["paySign"], err = c.Sign(params.Sort().ToURLParams(), params.Get("signType").String())
	delete(params, "appId")
	return
}
