package alipay

import (
	"bytes"
	gocrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"regexp"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

// func MapToStruct(dst interface{}, data map[string]interface{}) {
// 	//t := reflect.TypeOf(dst)
// 	v := reflect.ValueOf(dst)
// 	t := v.Type().Elem()
// 	for i := 0; i < t.NumField(); i++ {
// 		key := ""
// 		if tag, ok := t.Field(i).Tag.Lookup("json"); !ok || tag == "" {
// 			key = t.Field(i).Name
// 		} else {
// 			if index := strings.IndexAny(tag, ","); index <= 0 {
// 				key = tag
// 			} else {
// 				key = tag[:index]
// 			}
// 		}
// 		field := v.Elem().Field(i)
// 		if !field.CanSet() || !field.CanSet() {
// 			continue
// 		}
// 		if _, ok := data[key]; ok {
// 			if v, err := TypeConversion(fmt.Sprintf("%v", data[key]), field.Type().Name()); err == nil {
// 				field.Set(v)
// 			}
// 		}
// 	}
// }

// //类型转换
// func TypeConversion(value string, ntype string) (reflect.Value, error) {
// 	if ntype == "string" {
// 		return reflect.ValueOf(value), nil
// 	} else if ntype == "time.Time" {
// 		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
// 		return reflect.ValueOf(t), err
// 	} else if ntype == "Time" {
// 		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
// 		return reflect.ValueOf(t), err
// 	} else if ntype == "int" {
// 		i, err := strconv.Atoi(value)
// 		return reflect.ValueOf(i), err
// 	} else if ntype == "int8" {
// 		i, err := strconv.ParseInt(value, 10, 64)
// 		return reflect.ValueOf(int8(i)), err
// 	} else if ntype == "int32" {
// 		i, err := strconv.ParseInt(value, 10, 64)
// 		return reflect.ValueOf(int64(i)), err
// 	} else if ntype == "int64" {
// 		i, err := strconv.ParseInt(value, 10, 64)
// 		return reflect.ValueOf(i), err
// 	} else if ntype == "float32" {
// 		i, err := strconv.ParseFloat(value, 64)
// 		return reflect.ValueOf(float32(i)), err
// 	} else if ntype == "float64" {
// 		i, err := strconv.ParseFloat(value, 64)
// 		return reflect.ValueOf(i), err
// 	}

// 	//else if .......增加其他一些类型的转换

// 	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
// }

// func ToSortString(obj interface{}, ignoreEmpty, ignoreZero bool) (string, error) {
// 	t := reflect.TypeOf(obj)
// 	v := reflect.ValueOf(obj)
// 	var dict = make(map[string]string)
// 	var keys = make([]string, t.NumField())
// 	for i := 0; i < t.NumField(); i++ {
// 		key := ""
// 		if tag, ok := t.Field(i).Tag.Lookup("json"); !ok || tag == "" {
// 			continue
// 		} else {
// 			if index := strings.IndexAny(tag, ","); index <= 0 {
// 				key = tag
// 			} else {
// 				key = tag[:index]
// 			}
// 		}
// 		field := v.Field(i)

// 		val := ""
// 		switch field.Type().String() {
// 		case "string":
// 			val = field.String()
// 		case "uint":
// 			fallthrough
// 		case "uint64":
// 			fallthrough
// 		case "int":
// 			fallthrough
// 		case "int64":
// 			val = fmt.Sprintf("%d", field.Interface())
// 		case "*time.Time":
// 			{
// 				v := field.Interface()
// 				if t, ok := v.(*time.Time); ok && t != nil {
// 					val = t.Format("2006-01-02 15:04:05")
// 				}
// 			}
// 		case "time.Time":
// 			{
// 				val = field.Interface().(time.Time).Format("2006/01/02 15:04:05")
// 			}
// 		default:
// 			return "", fmt.Errorf("Not support filed type:%+v", field.Interface())
// 		}

// 		dict[key] = val
// 		keys[i] = key
// 	}
// 	sort.Strings(keys)
// 	var b strings.Builder
// 	first := true
// 	w := func(key string, val string) error {
// 		if ignoreEmpty && val == "" {
// 			return nil
// 		}
// 		if ignoreZero && val == "0" {
// 			return nil
// 		}
// 		if !first {
// 			if _, err := b.WriteString("&"); err != nil {
// 				return err
// 			}
// 		}
// 		if _, err := b.WriteString(key); err != nil {
// 			return err
// 		}
// 		if _, err := b.WriteString("="); err != nil {
// 			return err
// 		}
// 		if _, err := b.WriteString(val); err != nil {
// 			return err
// 		}
// 		first = false
// 		return nil
// 	}

// 	for _, key := range keys {
// 		val, _ := dict[key]
// 		if err := w(key, val); err != nil {
// 			return b.String(), err
// 		}

// 	}
// 	return b.String(), nil
// }

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
	//fmt.Println("=========")
	//fmt.Println(src)
	//fmt.Println("=========")
	//fmt.Println(sign)
	//fmt.Println("=========")
	hash := sha256.New()
	hash.Write([]byte(src))
	digest := hash.Sum(nil)

	// base64 decode,必须步骤，支付宝对返回的签名做过base64 encode必须要反过来decode才能通过验证
	data, _ := base64.StdEncoding.DecodeString(string(sign))

	// hex.EncodeToString(data)

	return rsa.VerifyPKCS1v15(pub, gocrypto.SHA256, digest, data)
}

func toGBK(src []byte) io.Reader {
	return transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewEncoder())
}

func toGBKData(src []byte) []byte {
	r := toGBK(src)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return data
}

func toUTF8(src []byte) io.Reader {
	return transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder())
}

func toUTF8Data(src []byte) []byte {
	r := toUTF8(src)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return data
}

// func MapToSortString(params map[string]string, ignoreEmpty, ignoreZero bool) string {
// 	var keys = make([]string, len(params))
// 	i := 0
// 	for key := range params {
// 		keys[i] = key
// 		i++
// 	}
// 	sort.Strings(keys)
// 	var b strings.Builder
// 	first := true
// 	w := func(key string, val string) error {
// 		if ignoreEmpty && val == "" {
// 			return nil
// 		}
// 		if ignoreZero && val == "0" {
// 			return nil
// 		}
// 		if !first {
// 			if _, err := b.WriteString("&"); err != nil {
// 				return err
// 			}
// 		}
// 		if _, err := b.WriteString(key); err != nil { //
// 			return err
// 		}
// 		if _, err := b.WriteString("="); err != nil {
// 			return err
// 		}
// 		if _, err := b.WriteString(val); err != nil {
// 			return err
// 		}
// 		first = false
// 		return nil
// 	}

// 	for _, key := range keys {
// 		if err := w(key, params[key]); err != nil {
// 			panic(err)
// 		}
// 	}
// 	return b.String()
// }

// func MapToSortJSON(params map[string]string, ignoreEmpty, ignoreZero bool) string {
// 	var keys = make([]string, len(params))
// 	i := 0
// 	for key := range params {
// 		keys[i] = key
// 		i++
// 	}
// 	sort.Strings(keys)
// 	var b strings.Builder
// 	first := true
// 	w := func(key string, val string) error {
// 		if ignoreEmpty && val == "" {
// 			return nil
// 		}
// 		if ignoreZero && val == "0" {
// 			return nil
// 		}
// 		if !first {
// 			if _, err := b.WriteString(","); err != nil {
// 				return err
// 			}
// 		}
// 		if _, err := b.WriteString("\"" + key + "\""); err != nil {
// 			return err
// 		}
// 		if _, err := b.WriteString(":"); err != nil {
// 			return err
// 		}
// 		if _, err := b.WriteString("\"" + val + "\""); err != nil {
// 			return err
// 		}
// 		first = false
// 		return nil
// 	}

// 	for _, key := range keys {
// 		if err := w(key, params[key]); err != nil {
// 			panic(err)
// 		}
// 	}
// 	return b.String()
// }
