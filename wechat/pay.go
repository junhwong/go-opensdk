package wechat

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/junhwong/go-opensdk/core"
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
func (c *WechatPayClient) UnifiedOrder(openID, outTradeNo, body, spbillIP string, totalFee int64) core.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/unifiedorder", core.Params{
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
func (c *WechatPayClient) OrderQuery(transactionID, outTradeNo string) core.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/orderquery", core.Params{
		"transaction_id": transactionID, // 微信的订单号，优先使用
		"out_trade_no":   outTradeNo,    // 商户系统内部订单号
	}).UseXML(true)
}

// CloseOrder 关单接口。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_3
func (c *WechatPayClient) CloseOrder(outTradeNo string) core.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/closeorder", core.Params{
		"out_trade_no": outTradeNo, // 商户系统内部订单号
	}).UseXML(true)
}

// Refund 申请退款。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_4
//
// 结果通知：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_16&index=10
func (c *WechatPayClient) Refund(transactionID, outTradeNo, outRefundNo string, totalFee, refundFee int64, refundDesc string) core.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/secapi/pay/refund", core.Params{
		"transaction_id":  transactionID, // 微信的订单号，优先使用
		"out_trade_no":    outTradeNo,    // 商户系统内部订单号
		"out_refund_no":   outRefundNo,   // 商户系统内部的退款单号
		"total_fee":       totalFee,      // 订单总金额
		"refund_fee":      refundFee,     // 退款总金额
		"refund_fee_type": "CNY",         // 货币类型
		"refund_desc":     refundDesc,    // 退款原因
	}).UseXML(true).UseTwowayAuthentication(true)
}

// RefundQuery 查询退款。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_5
func (c *WechatPayClient) RefundQuery(transactionID, outTradeNo, outRefundNo, refundID string) core.Executor {
	return c.buildPOST("https://api.mch.weixin.qq.com/pay/refundquery", core.Params{
		"transaction_id": transactionID, // 微信的订单号，优先使用
		"out_trade_no":   outTradeNo,    // 商户系统内部订单号
		"out_refund_no":  outRefundNo,   // 商户系统内部的退款单号
		"refund_id":      refundID,      // 退款单号
	}).UseXML(true)
}

//BuildMiniProgramRequestPaymentParams 构建小程序调起支付API参数。接口文档：https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=7_7&index=5
func (c *WechatPayClient) BuildMiniProgramRequestPaymentParams(prepayID string) (params core.Params, err error) {
	params = core.Params{
		"appId":     c.AppID,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  core.RandomString(10),
		"package":   "prepay_id=" + prepayID,
		"signType":  "MD5",
	}
	params["paySign"], err = c.Sign(params.Sort().ToURLParams(), params.Get("signType").String())
	delete(params, "appId")
	return
}

// MMPayTransfer 企业向微信用户个人付款。接口文档：https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_2
func (c *WechatPayClient) MMPayTransfer(partnerTradeNo, openID, reUserName, desc, spbillIP string, amount int64) core.Executor {
	checkName := "NO_CHECK"
	if reUserName != "" {
		checkName = "FORCE_CHECK"
	}
	return c.BuildExecutor("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", core.Params{
		"nonce_str":        core.RandomString(10),
		"mchid":            c.MchID,
		"mch_appid":        c.AppID,
		"partner_trade_no": partnerTradeNo,
		"openid":           openID, // 商户appid下，某用户的openid
		"check_name":       checkName,
		"re_user_name":     reUserName,
		"amount":           amount,
		"desc":             desc,
		"spbill_create_ip": spbillIP,
	}).UseXML(true).UseTwowayAuthentication(true)
}

// MMPayQuery 查询企业付款。接口文档：https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=14_3
func (c *WechatPayClient) MMPayQuery(partnerTradeNo string) core.Executor {
	return c.BuildExecutor("https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo", core.Params{
		"nonce_str":        core.RandomString(10),
		"mch_id":           c.MchID,
		"appid":            c.AppID,
		"partner_trade_no": partnerTradeNo,
	}).UseXML(true).UseTwowayAuthentication(true)
}
