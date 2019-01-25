package wechat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/junhwong/go-opensdk/opensdk"

	"github.com/junhwong/go-utils/crypto"
)

// MiniProgramLoginParams 小程序登录时的参数。
type MiniProgramLoginParams struct {
	Code          string `json:"code" form:"code"`
	Signature     string `json:"signature" form:"signature"`
	Iv            string `json:"iv" form:"iv"`
	EncryptedData string `json:"encryptedData" form:"encryptedData"`
}

// UserInfo 小程序加密用户信息。
type UserInfo struct {
	UserID      string `json:"userId"`
	OpenID      string `json:"openId"`
	UnionID     string `json:"unionid"`
	NickName    string `json:"nickName"`
	Gender      int    `json:"gender"`
	Language    string `json:"language"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	AvatarURL   string `json:"avatarUrl"`
	PhoneNumber string `json:"phoneNumber"`
}

// MiniProgramLoginSession code2Session应答结果
type MiniProgramLoginSession struct {
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errMsg"`
}

func decryptData(sessionKey, iv, encryptedData string) ([]byte, error) {
	key, _ := base64.StdEncoding.DecodeString(sessionKey)
	ivc, _ := base64.StdEncoding.DecodeString(iv)
	data, _ := base64.StdEncoding.DecodeString(encryptedData)
	return crypto.Decrypt(data, "AES-128-CBC", key, ivc)
	//return AesCBCDecrypt(data, key, ivc)
}

//JSCode2Session 文档 https://developers.weixin.qq.com/miniprogram/dev/api/open-api/login/code2Session.html
func (c *WechatClient) JSCode2Session(code string) (*MiniProgramLoginSession, error) {
	res, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + c.AppID + "&secret=" + c.Secret + "&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var sess MiniProgramLoginSession
	err = json.Unmarshal(body, &sess)

	return &sess, err
}

//MiniProgramLogin 小程序登录
func (c *WechatClient) MiniProgramLogin(params *MiniProgramLoginParams) (u *UserInfo, err error) {
	sess, err := c.JSCode2Session(params.Code)
	if err != nil {
		//panic(err)
		return
	}
	if sess.ErrCode != 0 {
		err = fmt.Errorf("%v", sess)
		return
	}
	s, err := decryptData(sess.SessionKey, params.Iv, params.EncryptedData)
	if err != nil {
		//panic(err)
		return
	}
	u = new(UserInfo)
	if err = json.Unmarshal(s, &u); err != nil || u.OpenID != sess.OpenID {
		if err == nil {
			err = fmt.Errorf("OpenID mismatched：%s,%s", u.OpenID, sess.OpenID)
		}
	}
	return
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Ticket       string `json:"ticket"`
	OpenID       string `json:"openid"`
	UnionID      string `json:"unionid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// CodeToAccessToken 文档 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
func (c *WechatClient) CodeToAccessToken(code string) (*AccessTokenResponse, error) {
	res, err := http.Get("https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + c.AppID + "&secret=" + c.Secret + "&code=" + code + "&grant_type=authorization_code")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var sess AccessTokenResponse
	err = json.Unmarshal(body, &sess)
	return &sess, err
}

// GetAccessToken 文档 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183
func (c *WechatClient) GetAccessToken() (*AccessTokenResponse, error) {
	res, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + c.AppID + "&secret=" + c.Secret)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var sess AccessTokenResponse
	err = json.Unmarshal(body, &sess)
	return &sess, err
}

// GetTicket 文档 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141115
func (c *WechatClient) GetTicket(token string) (*AccessTokenResponse, error) {
	res, err := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + token + "&type=jsapi")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var sess AccessTokenResponse
	err = json.Unmarshal(body, &sess)
	return &sess, err
}

// SendTemplateMessage 文档 https://developers.weixin.qq.com/miniprogram/dev/api/sendTemplateMessage.html
func (c *WechatClient) SendTemplateMessage(token, openID, templateID, formID, page, emphasisKeyword string, data opensdk.Params) ([]byte, error) {
	params := opensdk.Params{
		"access_token":     token,
		"touser":           openID,
		"template_id":      templateID,
		"page":             page,
		"form_id":          formID,
		"data":             data,
		"emphasis_keyword": emphasisKeyword,
	}
	requestBody := params.Sort(true).ToJSON(true)
	fmt.Println(requestBody)
	res, err := http.Post("https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token="+token, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
