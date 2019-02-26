package alipay

import "github.com/junhwong/go-opensdk/opensdk"

// TradePay 统一收单交易支付接口。接口文档： https://docs.open.alipay.com/api_1/alipay.trade.pay
//
// 注意：该方法默认为预授权转支付参数与官方API文档有差异。参考：https://docs.alipay.com/mini/introduce/pre-authorization
func (c *Client) TradePay(outTradeNo, authNo, subject string, amount string, storeID string, sellerID, buyerID string) opensdk.Executor {
	subject = string(opensdk.ToGBKData([]byte(subject))) // 支付宝内部使用 GBK编码
	return c.buildWithBizContent("alipay.trade.pay", opensdk.Params{
		"out_trade_no":      outTradeNo, // 新的交易流水号
		"auth_no":           authNo,     // 授权号
		"subject":           subject,    // 标题
		"total_amount":      amount,     // 支付金额
		"auth_confirm_mode": "COMPLETE", // 转交易支付完成结束预授权
		"store_id":          storeID,    // 需要与预授权的outStoreCode保持一致
		"product_code":      "PRE_AUTH_ONLINE",
		"seller_id":         sellerID,
		"buyer_id":          buyerID,
		"timeout_express":   "15d",

		//"trans_currency":    "CNY",
		//"scene":             "bar_code", // 场景, 默认为扫码支付

	})
}

// TradeQuery 统一收单交易查询。接口文档： https://docs.open.alipay.com/api_1/alipay.trade.query/
func (c *Client) TradeQuery(outTradeNo, tradeNo string) opensdk.Executor {
	return c.buildWithBizContent("alipay.trade.query", opensdk.Params{
		"out_trade_no": outTradeNo,
		"trade_no":     tradeNo,
	})
}

// TradeOrderInfoSync 支付宝订单信息同步接口。接口文档：https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync
// 特定使用文档：https://docs.alipay.com/mini/introduce/pre-authorization
// status: COMPLETE(用户已履约)、VIOLATED(用户已违约)；
func (c *Client) TradeOrderInfoSync(outOrderNo, outRequestNo string, status string) opensdk.Executor {
	return c.buildWithBizContent("alipay.trade.orderinfo.sync", opensdk.Params{
		"out_order_no":   outOrderNo,
		"out_request_no": outRequestNo,
		"biz_type":       "CREDIT_AUTH",
		"order_biz_info": `{"orderInfo":{\"status\":\"` + status + `\"}}`,
	})
}
