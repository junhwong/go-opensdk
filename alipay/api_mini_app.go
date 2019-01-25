package alipay

import "github.com/junhwong/go-opensdk/opensdk"

// MiniAppTemplateMessageSend 模板消息。接口文档：https://docs.alipay.com/mini/api/templatemessage
func (c *Client) MiniAppTemplateMessageSend(openID, templateID, formID, page, emphasisKeyword string, data string) opensdk.Executor {
	return c.buildWithBizContent("alipay.open.app.mini.templatemessage.send", opensdk.Params{
		"to_user_id":       openID,
		"form_id":          formID,
		"user_template_id": templateID,
		"page":             page,
		"data":             data,
	})
}
