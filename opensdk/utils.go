package opensdk

import (
	"bytes"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/crypto/pkcs12"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	rnd := strconv.FormatInt(rand.Int63(), n)
	if len(rnd) > n {
		rnd = rnd[:n]
	}
	return fmt.Sprintf("%s", rnd)
}

func ToGBK(src []byte) io.Reader {
	return transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewEncoder())
}

func ToGBKData(src []byte) []byte {
	r := ToGBK(src)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return data
}

func ToUTF8(src []byte) io.Reader {
	return transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder())
}

func ToUTF8Data(src []byte) ([]byte, error) {
	r := ToUTF8(src)
	data, err := ioutil.ReadAll(r)
	if err == io.EOF {
		err = nil
	}
	return data, err
}

// 将Pkcs12转成Pem
func PKCS12ToPem(p12 []byte, password string) tls.Certificate {

	blocks, err := pkcs12.ToPEM(p12, password)

	// 从恐慌恢复
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()

	if err != nil {
		panic(err)
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}
	return cert
}
