package opensdk

import (
	utils "github.com/junhwong/go-utils"
)

// Results 执行接口返回的响应结果。
type Results interface {
	// Set 重新设置结果中的某个值，以用于数据绑定。
	// Set(key string, v interface{}, err error) Results
	// Get 获取结果中的某个值。
	Get(key string) utils.Converter
	// Bind 绑定到一个结构体。
	// applyType 绑定是字段匹配的类型, 值：default、json、xml。
	// Bind(v interface{}, applyType ...string) error
	// Sign() string
	Code() string
	Message() string
	SubCode() string
	SubMessage() string
	Body() []byte
	Error() error
	Success() bool
	Values() Params
}

type DefaultResults struct {
	Params
	Err           error
	Data          []byte
	ResultCode    string
	ResultMsg     string
	ResultSubCode string
	ResultSubMsg  string
	ResultSuccess bool
}

func (r *DefaultResults) Values() Params {
	return r.Params
}

// func (r *DefaultResults) Set(key string, v interface{}, err error) Results {
// 	return r
// }
func (r *DefaultResults) Get(key string) utils.Converter {
	return utils.Convert(r.Params[key], r.Err, true)
}

// func (r *DefaultResults) Bind(v interface{}, applyType ...string) error {
// 	if r.Error() != nil {
// 		return r.Error()
// 	}
// 	// var data []byte
// 	// if r.Success(){
// 	// 	data=[]byte(r.ResultString)
// 	// }else {
// 	// 	data=[]byte(r.ResultString)
// 	// }
// 	data := []byte(r.ResultString)
// 	return json.Unmarshal(data, v)
// }
func (r *DefaultResults) Error() error {
	return r.Err
}
func (r *DefaultResults) Body() []byte {
	return r.Data
}
func (r *DefaultResults) Success() bool {
	return r.Err == nil && r.ResultSuccess
}

func (r *DefaultResults) Code() string {
	return r.ResultCode
}
func (r *DefaultResults) Message() string {
	return r.ResultMsg
}
func (r *DefaultResults) SubCode() string {
	return r.ResultSubCode
}
func (r *DefaultResults) SubMessage() string {
	return r.ResultSubMsg
}

// func (r *DefaultResults) Success() bool {
// 	return r.Error() == nil && r.Code() == "10000"
// }

// func (r *DefaultResults) Code() string {
// 	return r.Get("code").String()
// }
// func (r *DefaultResults) Message() string {
// 	return r.Get("msg").String()
// }
// func (r *DefaultResults) SubCode() string {
// 	return r.Get("sub_code").String()
// }
// func (r *DefaultResults) SubMessage() string {
// 	return r.Get("sub_msg").String()
// }

// func (r *DefaultResults) Sign() string {
// 	return r.SignString
// }

// Default 执行接口返回的响应。
// type Response struct {
// 	Sign       string     `json:"sign"`
// 	Code       string     `json:"code"`
// 	Message    string     `json:"msg"`
// 	SubCode    string     `json:"sub_code"`
// 	SubMessage string     `json:"sub_msg"`
// 	Err        error      `json:"-"`
// 	Body       []byte     `json:"-"`
// 	Results    Parameters `json:"-"`
// }

// func (r *Response) Set(key string, v interface{}, err error) Results {
// 	return r
// }
// func (r *Response) Get(key string) utils.Converter {
// 	return nil
// }
// func (r *Response) Bind(v interface{}, applyType ...string) error {
// 	return nil
// }
// func (r *Response) Error() error {
// 	return r.Err
// }
// func (r *Response) Data() []byte {
// 	return r.Body
// }
// func (r *Response) Success() bool {
// 	return false
// }

// func (r *Response) Code() string {
// 	return ""
// }

// type ResponseParser func(*http.Response) *Response
// type ResponseParseFn func(*http.Response) *Response
// type Extract func(Parameters) (code, msg, subCode, subMsg, sign string, result Parameters, isSuccess bool)

// func DefaultResponseParser(r *http.Response, e Extract) *Response {
// 	defer r.Body.Close()
// 	body, err := ioutil.ReadAll(r.Body)
// 	ret := Response{
// 		Err:     err,
// 		Body:    body,
// 		Results: Parameters{},
// 	}
// 	if ret.Err == nil {
// 		ret.Err = json.Unmarshal(body, &ret.Results)
// 		if ret.Err == nil {

// 		}
// 	}
// 	utils.Convert(nil, ret.Err)
// 	return &ret
// }
