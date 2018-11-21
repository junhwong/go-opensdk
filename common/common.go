package common

import (
	"encoding/json"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/junhwong/go-utils"
)

type Client interface {
	Build(methodOrURL string) Executor
}

type Params map[string]interface{}

func (p Params) Get(key string) utils.Converter {
	return utils.Convert(p[key], nil)
}

func (p Params) getKeys() []string {
	var keys = []string{}
	for k, v := range p {
		if k == "" || v == nil {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (p Params) Sort() Pairs {
	arr := Pairs{}
	for _, key := range p.getKeys() {
		v := p[key]
		val := ""
		if m, ok := v.(Params); ok {
			val = m.Sort().ToJSON()
		} else if m, ok := v.(map[string]interface{}); ok {
			val = Params(m).Sort().ToJSON()
		} else {
			if d, err := json.Marshal(v); err != nil {
				val = "<MarshalError>"
			} else {
				val = string(d)
			}
			if strings.HasPrefix(val, `{`) && strings.HasSuffix(val, `}`) {
				//结构JSON字符串
				sub := Params{}
				if err := json.Unmarshal([]byte(val), &sub); err == nil {
					val = sub.Sort().ToJSON()
				}
			}
		}
		if val != "" {
			arr = append(arr, [2]string{key, val})
		}
	}
	return arr
}

type Pairs [][2]string

func (p Pairs) ToURLParams(urlencode ...bool) string {
	e := false
	if len(urlencode) > 0 {
		e = urlencode[0]
	}
	arr := []string{}
	for _, it := range p {
		val := it[1]
		if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
			val = strings.Trim(val, `"`)
		}
		if e {
			val = url.QueryEscape(val)
		}
		arr = append(arr, it[0]+"="+val)
	}
	return strings.Join(arr, "&")
}

func (p Pairs) ToJSON() string {
	arr := []string{}
	for _, it := range p {
		val := it[1]
		if !(strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`)) {
			val = `"` + it[1] + `"`
		}
		arr = append(arr, `"`+it[0]+`":`+val)
	}
	return `{` + strings.Join(arr, ",") + `}`
}

func (p *Pairs) Append(k, v string) {
	r := append(*p, [2]string{k, v})
	*p = *(&r)
}

type JsonTime time.Time

//MarshalJSON 实现它的json序列化方法
func (t JsonTime) MarshalJSON() ([]byte, error) {
	s := time.Time(t).Format("2006-01-02 15:04:05")
	return []byte("\"" + s + "\""), nil
}
