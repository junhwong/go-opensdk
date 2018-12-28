package wechat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/junhwong/go-utils/crypto"
)

// MiniProgramLoginParams 小程序登录时的参数。
type MiniProgramLoginParams struct {
	Code          string `json:"code"`
	Signature     string `json:"signature"`
	Iv            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
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
