package alipay

import (
	"strings"

	"github.com/junhwong/go-opensdk/opensdk"
)

// FundAuthOrderAppFreezeRequest 构建app端发起资金预授权支付的参数。接口文档： https://docs.open.alipay.com/api_28/alipay.fund.auth.order.app.freeze
//
// 小程序端发起文档：https://docs.alipay.com/mini/introduce/pre-authorization
//
// 预授权类目：https://docs.open.alipay.com/10719
func (c *Client) FundAuthOrderAppFreezeRequest(outOrderNo, outRequestNo, orderTitle string, amount string, payeeUserID string, extraParam opensdk.Params, notifyURL string) (params opensdk.Params, err error) {

	params = opensdk.Params{

		"out_order_no":        outOrderNo,
		"out_request_no":      outRequestNo,
		"order_title":         orderTitle,
		"amount":              amount,
		"product_code":        "PRE_AUTH_ONLINE",
		"payee_user_id":       payeeUserID,
		"pay_timeout":         "15d",                                                      // 超时 取值范围：1m～15d。m-分钟，h-小时，d-天。 该参数数值不接受小数点， 如 1.5h，可转换为90m ，如果为空，默认15m
		"extra_param":         strings.Replace(extraParam.Sort().ToJSON(), `"`, `\"`, -1), // TODO: 内嵌JSON字符串处理
		"enable_pay_channels": `[{\"payChannelType\":\"CREDITZHIMA\"},{\"payChannelType\":\"PCREDIT_PAY\"},{\"payChannelType\":\"MONEY_FUND\"}]`,
	}
	header := opensdk.Params{"notify_url": notifyURL}
	c.fillParams("alipay.fund.auth.order.app.freeze", header)
	header["biz_content"] = params.Sort().ToJSON()
	header["sign"], err = c.Sign(header.Sort().ToURLParams(), header.Get("sign_type").String())
	return header, err
}

// FundAuthOrderUnfreeze 资金授权解冻接口。接口文档： https://docs.open.alipay.com/api_28/alipay.fund.auth.order.unfreeze
func (c *Client) FundAuthOrderUnfreeze(outRequestNo, authNo, remark string, amount string) opensdk.Executor {
	remark = string(opensdk.ToGBKData([]byte(remark))) // 支付宝内部使用 GBK编码
	return c.buildWithBizContent("alipay.fund.auth.order.unfreeze", opensdk.Params{
		"out_request_no": outRequestNo,
		"auth_no":        authNo,
		"remark":         remark,
		"amount":         amount,                                                         // 解冻的金额
		"extra_param":    `{\"unfreezeBizInfo\":\"{\\\"bizComplete\\\":\\\"true\\\"}\"}`, // 是否履约
	})
}

// FundAuthOperationCancel 资金授权撤销接口。接口文档： https://docs.open.alipay.com/api_28/alipay.fund.auth.operation.cancel
func (c *Client) FundAuthOperationCancel(outOrderNo, outRequestNo string, remark string) opensdk.Executor {
	remark = string(opensdk.ToGBKData([]byte(remark))) // 支付宝内部使用 GBK编码
	return c.buildWithBizContent("alipay.fund.auth.operation.cancel", opensdk.Params{
		"out_order_no":   outOrderNo,
		"out_request_no": outRequestNo,
		"remark":         remark,
	})
}

// FundAuthOperationDetailQuery 资金授权操作查询接口。接口文档：https://docs.open.alipay.com/api_28/alipay.fund.auth.operation.detail.query
func (c *Client) FundAuthOperationDetailQuery(outOrderNo, outRequestNo string, authNo, operationID string) opensdk.Executor {
	return c.buildWithBizContent("alipay.fund.auth.operation.detail.query", opensdk.Params{
		"out_order_no":   outOrderNo,
		"out_request_no": outRequestNo,
		"auth_no":        authNo,
		"operation_id":   operationID,
	})
}
