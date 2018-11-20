package wechat

import (
	"crypto/rsa"
	"net/url"
	"strings"

	"github.com/junhwong/go-opensdk/common"
)

type Client struct {
	AppID      string
	Secret     string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	Gateway    string
}

func NewClient(appid, secret string) *Client {
	return &Client{
		AppID:  appid,
		Secret: secret,
	}
}

func (c *Client) Build(api string) common.Executor {
	if api == "" {

	}
	if strings.IndexAny(api, "http://") != -1 || strings.IndexAny(api, "https://") != -1 {

	} else {

	}
	url, err := url.Parse(api)
	if err != nil {
		panic(err)
	}
	return &common.DefaultExecutor{
		RequestURL:      url,
		RequestMethod:   "POST",
		RequestEncoding: "UTF-8",
		Params:          common.Parameters{},
	}
}
