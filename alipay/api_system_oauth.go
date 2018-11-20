package alipay

//SystemOauthToken 换取授权访问令牌。接口文档：https://docs.open.alipay.com/api_9/alipay.system.oauth.token
func (c *Client) SystemOauthToken(code string, refreshCode bool) *Executor {

	params := BuildParams("alipay.system.oauth.token")

	if refreshCode {
		params["grant_type"] = "refresh_token"
		params["refresh_token"] = code
	} else {
		params["grant_type"] = "authorization_code"
		params["code"] = code
	}

	return c.Execute(params, "alipay_system_oauth_token_response")

}
