package opensdk

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

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
