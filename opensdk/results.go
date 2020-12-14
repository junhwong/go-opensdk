package opensdk

import (
	utils "github.com/junhwong/go-utils"
)

// Results 执行接口返回的响应结果。
type Results interface {
	Get(key string) utils.Converter
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
	Data          []byte `json:"-"`
	StatusCode    int
	Header        interface{}
	ResultCode    string
	ResultMsg     string
	ResultSubCode string
	ResultSubMsg  string
	ResultSuccess bool
}

func (r *DefaultResults) Values() Params {
	return r.Params
}

func (r *DefaultResults) Get(key string) utils.Converter {
	return utils.Convert(r.Params[key], r.Err, true)
}

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
