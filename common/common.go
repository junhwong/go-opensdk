package common

import "github.com/junhwong/go-utils"

type Parameters map[string]interface{}

func (p Parameters) Set(key string, v interface{}) {
	p[key] = v
}

func (p Parameters) Get(key string) utils.Converter {
	return utils.Convert(p[key], nil)
}
