package wechat

import (
	"crypto/rsa"

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
	return nil
}
