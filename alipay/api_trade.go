package alipay

import "github.com/junhwong/go-opensdk/core"

// TradePay 统一收单交易支付接口。接口文档： https://docs.open.alipay.com/api_1/alipay.trade.pay
//
// 注意：该方法默认为预授权转支付参数与官方API文档有差异。参考：https://docs.alipay.com/mini/introduce/pre-authorization
func (c *Client) TradePay(outTradeNo, authNo, subject string, amount string, storeID string, sellerID, buyerID string) core.Executor {
	subject = string(core.ToGBKData([]byte(subject))) // 支付宝内部使用 GBK编码
	return c.buildWithBizContent("alipay.trade.pay", core.Params{
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
func (c *Client) TradeQuery(outTradeNo, tradeNo string) core.Executor {
	return c.buildWithBizContent("alipay.trade.query", core.Params{
		"out_trade_no": outTradeNo,
		"trade_no":     tradeNo,
	})
}

// TradeOrderInfoSync 支付宝订单信息同步接口。接口文档：https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync
// 特定使用文档：https://docs.alipay.com/mini/introduce/pre-authorization
// status: COMPLETE(用户已履约)、VIOLATED(用户已违约)；
func (c *Client) TradeOrderInfoSync(outOrderNo, tradeNo string, status string) opensdk.Executor {
	return c.buildWithBizContent("alipay.trade.orderinfo.sync", opensdk.Params{
		"trade_no":       tradeNo,
		"out_request_no": outOrderNo,
		"biz_type":       "CREDIT_AUTH",
		// "order_biz_info": `{"orderInfo":{\"status\":\"` + status + `\"}}`,
		"order_biz_info": `{\"status\":\"` + status + `\"}`,
	})
}

// TradeCreate 商户通过该接口进行交易的创建下单。接口文档： https://docs.open.alipay.com/api_1/alipay.trade.create/
//
func (c *Client) TradeCreate(outTradeNo, subject, buyer_id string, amount string) opensdk.Executor {
	subject = string(opensdk.ToGBKData([]byte(subject))) // 支付宝内部使用 GBK编码
	return c.buildWithBizContent("alipay.trade.create", opensdk.Params{
		"out_trade_no":    outTradeNo, // 新的交易流水号
		"subject":         subject,    // 标题
		"total_amount":    amount,     // 支付金额
		"buyer_id":        buyer_id,
		"timeout_express": "10m",
	})
}

// TradeClose 统一收单交易关闭接口。接口文档： https://docs.open.alipay.com/api_1/alipay.trade.close/
//
func (c *Client) TradeClose(outTradeNo, tradeNo string) opensdk.Executor {
	// subject = string(opensdk.ToGBKData([]byte(subject))) // 支付宝内部使用 GBK编码
	return c.buildWithBizContent("alipay.trade.close", opensdk.Params{
		"out_trade_no": outTradeNo, // 新的交易流水号
		"trade_no":     tradeNo,    // 标题
	})
}

// TradeRefund 统一收单交易退款接口。接口文档： https://docs.open.alipay.com/api_1/alipay.trade.refund/
//
func (c *Client) TradeRefund(outTradeNo, tradeNo, outRequestNo string, amount string) opensdk.Executor {
	// subject = string(opensdk.ToGBKData([]byte(subject))) // 支付宝内部使用 GBK编码
	return c.buildWithBizContent("alipay.trade.refund", opensdk.Params{
		"out_trade_no":   outTradeNo,   // 新的交易流水号
		"trade_no":       tradeNo,      // 标题
		"out_request_no": outRequestNo, // 标题
		"refund_amount":  amount,       // 标题
	})
}

// TradeRefundQuery 统一收单交易退款查询。接口文档： https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query/
func (c *Client) TradeRefundQuery(outTradeNo, tradeNo, outRequestNo string) opensdk.Executor {
	return c.buildWithBizContent("alipay.trade.fastpay.refund.query", opensdk.Params{
		"out_trade_no":   outTradeNo,
		"trade_no":       tradeNo,
		"out_request_no": outRequestNo, // 标题
	})
}
