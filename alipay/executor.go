package alipay

// type Executor struct {
// 	opensdk.BaseExecutor
// 	// params2 opensdk.Params
// 	paramsStr string
// 	// bizContent   map[string]string
// 	// executed     bool
// 	// err          error
// 	client *Client
// 	// body         []byte
// 	// resultsField string
// }

// func (e *Executor) Client() opensdk.Client {
// 	return nil
// }
// func (e *Executor) UseTLS(b bool) opensdk.Executor {
// 	return e
// }

// func (e *Executor) UseXML(b bool) opensdk.Executor {
// 	return e
// }

// func (e *Executor) Execute(verifySign ...bool) opensdk.Results {
// 	delete(e.Params, "sign")
// 	ret := &ResponseResults{
// 		DefaultResults: opensdk.DefaultResults{
// 			Params: opensdk.Params{},
// 		},
// 	}
// 	params := e.Params.Sort()
// 	sign, err := sha256WithRSA(params.ToURLParams(false), e.client.PrivateKey)
// 	if err != nil {
// 		ret.Err = err
// 		return ret
// 	}
// 	params = append(params, [2]string{"sign", sign})
// 	u, _ := url.Parse(e.client.Gateway)
// 	e.paramsStr = params.ToURLParams(true)
// 	fmt.Println(e.paramsStr)
// 	ret.Data, _, ret.Err = e.DoPost(u, strings.NewReader(params.ToURLParams(true)))
// 	if ret.Err != nil {
// 		return ret
// 	}
// 	ret.ResultString, ret.SignString = extract(ret.Data)
// 	//签名验证
// 	if len(verifySign) > 0 && verifySign[0] {
// 		ret.Err = verifyRSA2(ret.ResultString, ret.SignString, e.client.PublicKey)
// 		if ret.Err != nil {
// 			return ret
// 		}
// 	}
// 	reader := toUTF8([]byte(ret.ResultString)) // 支付宝返回编码是GBK，不管传递参数是不是GBK。这是BUG?
// 	data, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		ret.Err = err
// 		return ret
// 	}
// 	ret.Err = json.Unmarshal(data, &ret.Params)
// 	return ret
// }

// type ResponseResults struct {
// 	opensdk.DefaultResults
// }

// func (r *ResponseResults) Success() bool {
// 	return r.Error() == nil && r.Code() == "10000"
// }
