package alipay

import "github.com/junhwong/go-opensdk/common"

type UserInfo struct {
	OpenID             string `json:"user_id"`
	AvatarURL          string `json:"avatar"`
	Province           string `json:"province"`
	City               string `json:"city"`
	NickName           string `json:"nick_name"`
	IsStudentCertified string `json:"is_student_certified"`
	UserType           string `json:"user_type"`
	UserStatus         string `json:"user_status"`
	IsCertified        string `json:"is_certified"`
	Gender             string `json:"gender"`
}

//SystemOauthToken 换取授权访问令牌。接口文档：https://docs.open.alipay.com/api_2/alipay.user.info.share
func (c *Client) UserInfoShare(authToken string) common.Executor {
	return c.Build("alipay.user.info.share", common.Params{
		"auth_token": authToken,
	})
}

// func (c *Client) UserInfoShare(authToken string) *Executor {

// 	params := BuildParams("alipay.user.info.share")

// 	params["auth_token"] = authToken

// 	return c.Execute(params, "alipay_user_info_share_response")

// }
