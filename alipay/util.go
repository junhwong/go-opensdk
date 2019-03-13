package alipay

import (
	gocrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"regexp"
)

var regResult = regexp.MustCompile(`(^\{\"[a-z|_]+\":)|(,\"sign\":\"[a-zA-Z0-9|\+|\/|\=]+\"\}$)`)
var regSign = regexp.MustCompile(`\"sign\":\"([a-zA-Z0-9|\+|\/|\=]+)\"`)

func extract(data []byte) (result, sign string) {
	s := string(data)
	result = regResult.ReplaceAllString(s, "")
	m := regSign.FindStringSubmatch(s)
	if len(m) > 1 {
		sign = m[1]
	}
	return
}

func sha256WithRSA(origData string, key *rsa.PrivateKey) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(origData))
	digest := hash.Sum(nil)
	s, err := rsa.SignPKCS1v15(rand.Reader, key, gocrypto.SHA256, digest)
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return string(data), nil
}

func verifyRSA2(src, sign string, pub *rsa.PublicKey) error {
	hash := sha256.New()
	hash.Write([]byte(src))
	digest := hash.Sum(nil)

	// base64 decode,必须步骤，支付宝对返回的签名做过base64 encode必须要反过来decode才能通过验证
	data, _ := base64.StdEncoding.DecodeString(string(sign))

	// hex.EncodeToString(data)

	return rsa.VerifyPKCS1v15(pub, gocrypto.SHA256, digest, data)
}
