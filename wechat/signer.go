package wechat

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ParsePrivateKey(pemData []byte) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode(pemData) //rest
	if block == nil || block.Type != "PRIVATE KEY" {
		if block != nil {
			return nil, fmt.Errorf("failed to decode PEM block: %s", block.Type)
		}
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pri.(*rsa.PrivateKey), nil
	// rsa.SignPKCS1v15()
	// rsa.SignPSS()Sha256RSA
	//fmt.Printf("Got a %T, with remaining data: %q", pri, rest)
}

type Signer interface {
	Sign(message []byte) Signature
}
type Signature interface {
	Signature() string
	CertificateSerialNumber() string
}
