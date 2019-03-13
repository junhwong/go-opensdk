package core

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

func Sha256Hmac(s string, k []byte) (string, error) {
	h := hmac.New(sha256.New, k)
	_, err := io.WriteString(h, s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", h.Sum(nil)), nil
}

func Sha256RSA(origData string, key *rsa.PrivateKey) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(origData))
	digest := hash.Sum(nil)
	s, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return string(data), nil
}

func Sha1(s string) (string, error) {
	h := sha1.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func MD5(s string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", h.Sum(nil)), nil
}
