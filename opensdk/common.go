package opensdk

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/junhwong/go-utils"
)

type Client interface {
	Sign(params, signType string) (string, error)
	VerifySign(params, signType string) (bool, error)
	// ff
	HttpClient(twowayAuthentication bool) (*http.Client, error)
}

type ClientBase struct {
	GetHttpClient func(twowayAuthentication bool) (*http.Client, error)
	hc            *http.Client
}

func (c *ClientBase) HttpClient(twowayAuthentication bool) (*http.Client, error) {
	if c.GetHttpClient != nil {
		return c.GetHttpClient(twowayAuthentication)
	} else {
		return http.DefaultClient, nil
	}
	// FIXME: 复用会出EOF BUG
	var err error
	if c.hc == nil {
		if c.GetHttpClient != nil {
			c.hc, err = c.GetHttpClient(twowayAuthentication)
		} else {
			return &http.Client{}, nil
		}
	} else {
		if c.hc.Transport != nil {
			ts := c.hc.Transport
			if x, ok := ts.(*http.Transport); ok {
				transport := &http.Transport{
					TLSClientConfig:    x.TLSClientConfig.Clone(),
					DisableCompression: true,
					DisableKeepAlives:  true,
				}
				return &http.Client{Transport: transport}, nil
			}
		}
	}
	return c.hc, err
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

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m *Params) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Params{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func (p Params) Sort(isNesting ...bool) Pairs {
	arr := Pairs{}
	for _, key := range p.getKeys() {
		v := p[key]
		val := ""
		if m, ok := v.(Params); ok {
			val = m.Sort(isNesting...).ToJSON(isNesting...)
		} else if m, ok := v.(map[string]interface{}); ok {
			val = Params(m).Sort(isNesting...).ToJSON(isNesting...)
		} else {
			if s, ok := v.(string); ok {
				// if strings.HasPrefix(s, `{`) && strings.HasSuffix(s, `}`) {
				// 	//结构JSON字符串,重新排序
				// 	sub := Params{}
				// 	if err := json.Unmarshal([]byte(s), &sub); err == nil {
				// 		val = sub.Sort().ToJSON()
				// 	}
				// } else {
				// 	val = s
				// }
				val = s
			} else {
				// val = fmt.Sprint(v)
				if d, err := json.Marshal(v); err != nil {
					val = "<MarshalError>"
				} else {
					val = string(d)
				}
			}
			if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
				val = strings.Trim(val, `"`)
			}
		}
		if val != "" {
			arr = append(arr, [2]string{key, val})
		}
	}
	return arr
}
func (p Params) Sort2(isNesting ...bool) Pairs {
	arr := Pairs{}
	for _, key := range p.getKeys() {
		v := p[key]
		val := ""
		if m, ok := v.(Params); ok {
			val = m.Sort2(isNesting...).ToJSON2(isNesting...)
		} else if m, ok := v.(map[string]interface{}); ok {
			val = Params(m).Sort2(isNesting...).ToJSON2(isNesting...)
		} else {
			if d, err := json.Marshal(v); err != nil {
				val = "<MarshalError>"
			} else {
				val = string(d)
			}
			// if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
			// 	val = strings.Trim(val, `"`)
			// }
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
		// if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
		// 	val = strings.Trim(val, `"`)
		// }
		if e {
			val = url.QueryEscape(val)
		}
		arr = append(arr, it[0]+"="+val)
	}
	return strings.Join(arr, "&")
}

func (p Pairs) ToJSON(isNesting ...bool) string {
	arr := []string{}
	for _, it := range p {
		// val := it[1]
		// if !(strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`)) {
		// 	val = `"` + it[1] + `"`
		// }
		if (len(isNesting) > 0 && isNesting[0]) && (strings.HasPrefix(it[1], `{`) && strings.HasSuffix(it[1], `}`)) {
			arr = append(arr, `"`+it[0]+`":`+it[1])
		} else {
			arr = append(arr, `"`+it[0]+`":`+`"`+it[1]+`"`)
		}

	}
	return `{` + strings.Join(arr, ",") + `}`
}
func (p Pairs) ToJSON2(isNesting ...bool) string {
	arr := []string{}
	for _, it := range p {
		arr = append(arr, `"`+it[0]+`":`+it[1])
		// val := it[1]
		// if !(strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`)) {
		// 	val = `"` + it[1] + `"`
		// }
		// if (len(isNesting) > 0 && isNesting[0]) && (strings.HasPrefix(it[1], `{`) && strings.HasSuffix(it[1], `}`)) {
		// 	arr = append(arr, `"`+it[0]+`":`+it[1])
		// } else {
		// 	arr = append(arr, `"`+it[0]+`":`+`"`+it[1]+`"`)
		// }

	}
	return `{` + strings.Join(arr, ",") + `}`
}

func (p Pairs) ToXML() string {
	arr := []string{}
	for _, it := range p {
		val := it[1]
		if strings.IndexAny(val, `<`) != -1 || strings.IndexAny(val, `>`) != -1 {
			val = `<![CDATA[` + it[1] + `]]>`
		}
		arr = append(arr, `<`+it[0]+`>`+val+`</`+it[0]+`>`)
	}
	return `<xml>` + strings.Join(arr, "") + `</xml>`
}

func (p *Pairs) Append(k, v string) {
	r := append(*p, [2]string{k, v})
	*p = *(&r)
}

func (p Pairs) ToHeaderParams(urlencode ...bool) string {
	e := false
	if len(urlencode) > 0 {
		e = urlencode[0]
	}
	arr := []string{}
	for _, it := range p {
		val := it[1]
		// if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
		// 	val = strings.Trim(val, `"`)
		// }
		if e {
			val = url.QueryEscape(val)
		}
		arr = append(arr, it[0]+`="`+val+`"`)
	}
	return strings.Join(arr, ",")
}

type JsonTime time.Time

//MarshalJSON 实现它的json序列化方法
func (t JsonTime) MarshalJSON() ([]byte, error) {
	s := time.Time(t).Format("2006-01-02 15:04:05")
	return []byte("\"" + s + "\""), nil
}
