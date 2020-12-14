package wechat

import (
	"context"
	"fmt"

	"time"

	"github.com/junhwong/go-logs"
	"github.com/junhwong/go-opensdk/opensdk"
	"github.com/junhwong/go-utils"
)

// 请求参数，必填
type PayScoreRentParams struct {
	OutOrderNo      string    `json:"out_order_no" xml:"out_order_no"`         // 商户服务订单号
	GoodsName       string    `json:"goods_name" xml:"goods_name"`             // 商品名称
	StartTime       time.Time `json:"start_time" xml:"start_time"`             // 租用时间
	EndTime         time.Time `json:"end_time" xml:"end_time"`                 // 预定归还时间
	ServiceLocation string    `json:"service_location" xml:"service_location"` // 租用地点
	DepositAmount   int64     `json:"deposit_amount" xml:"deposit_amount"`     // 押金总额
	RentUnitFee     int64     `json:"rent_unit_fee" xml:"rent_unit_fee"`       // 租金规则 计费单价费用 分/小时
	RentFeeDesc     string    `json:"rent_fee_desc" xml:"rent_fee_desc"`       // 租金规则 计费单价费用 分/小时
	NotifyURL       string    `json:"-" xml:"-"`
}

// 请求参数，必填
type PayScoreRentCreateResults struct {
	FinishTicket        string         `json:"finish_ticket,omitempty"`                         // 完结凭证
	FinishTransactionID string         `json:"finish_transaction_id,omitempty"`                 // 结单交易单号
	ServiceID           string         `json:"service_id,omitempty"`                            // 支付渠道
	OutOrderNo          string         `json:"out_order_no" xml:"out_order_no"`                 // 商户服务订单号
	OrderID             string         `json:"order_id" xml:"order_id"`                         // 微信支付服务订单号
	MiniprogramAppid    string         `json:"miniprogram_appid" xml:"miniprogram_appid"`       // 小程序跳转appid
	MiniprogramPath     string         `json:"miniprogram_path" xml:"miniprogram_path"`         // 小程序跳转路径
	MiniprogramUsername string         `json:"miniprogram_username" xml:"miniprogram_username"` // 小程序跳转username
	Package             string         `json:"package" xml:"package"`                           // 跳转微信侧小程序订单数据
	MchID               string         `json:"mch_id" xml:"mch_id"`                             // 跳转微信侧小程序订单数据
	State               string         `json:"state" xml:"state"`
	Params              opensdk.Params `json:"-" xml:"-" gorm:"-"`
	opensdk.Results     `json:"-" xml:"-" gorm:"-"`
}

// BuildMiniaConfirmParams 构造订单确认参数(小程序端)。接口文档：https://pay.weixin.qq.com/wiki/doc/apiv3/payscore.php?chapter=30_2&index=2
//
// envVersion: release(默认), trial, develop
func (r *PayScoreRentCreateResults) BuildMiniaConfirmParams(c *WechatPayClient, envVersion ...string) (p opensdk.Params, err error) {
	p = opensdk.Params{
		"businessType": "wxpayScoreUse",
		"envVersion":   "release",
	}

	for _, it := range envVersion {
		switch it {
		case "trial":
			p["envVersion"] = "trial"
		case "develop":
			p["envVersion"] = "develop"
		}
	}
	data := opensdk.Params{
		"mch_id":    c.MchID,
		"package":   r.Package,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonce_str": opensdk.RandomString(10),
		"sign_type": "HMAC-SHA256",
	}
	temp := data.Sort().ToURLParams() + "&key=" + c.MchKey
	data["sign"], err = opensdk.Sha256Hmac(temp, []byte(c.MchKey))
	p["extraData"] = data

	return
}

// BuildMiniaEnableParams 构造订单开启服务参数(小程序端)。接口文档：https://pay.weixin.qq.com/wiki/doc/apiv3/payscore.php?chapter=30_2&index=2
//
// envVersion: release(默认), trial, develop
func (r *PayScoreRentCreateResults) BuildMiniaEnableParams(c *WechatPayClient, envVersion ...string) (p opensdk.Params, err error) {
	p = opensdk.Params{
		"businessType": "wxpayScoreEnable",
		"envVersion":   "release",
	}

	for _, it := range envVersion {
		switch it {
		case "trial":
			p["envVersion"] = "trial"
		case "develop":
			p["envVersion"] = "develop"
		}
	}
	data := opensdk.Params{
		"mch_id":         r.MchID,
		"service_id":     r.ServiceID,
		"out_request_no": r.OutOrderNo,
		"timestamp":      fmt.Sprintf("%d", time.Now().Unix()),
		"nonce_str":      opensdk.RandomString(10),
		"sign_type":      "HMAC-SHA256",
	}
	temp := data.Sort().ToURLParams() + "&key=" + c.MchKey
	data["sign"], err = opensdk.Sha256Hmac(temp, []byte(c.MchKey))
	// fmt.Println(temp)
	// fmt.Println(c.MchKey)
	// fmt.Println(data["sign"])
	p["extraData"] = data

	return
}

type ParamsOption = func(opensdk.Params)

// PayScoreRentCreate 创建租借订单。接口文档：https://pay.weixin.qq.com/wiki/doc/apiv3/payscore.php?chapter=18_1&index=2
func (c *WechatPayClient) PayScoreRentCreate(ctx context.Context, params PayScoreRentParams, options ...ParamsOption) (results *PayScoreRentCreateResults, err error) {
	// 签名方式：HMAC-SHA256
	p := opensdk.Params{
		"version":      "1.0",
		"sign_type":    "HMAC-SHA256",
		"nonce_str":    opensdk.RandomString(10),
		"mch_id":       c.MchID,
		"appid":        c.AppID,
		"service_id":   c.ServiceID,
		"rent_unit":    "FEN_1_HOUR",
		"out_order_no": params.OutOrderNo,
		"goods_name":   params.GoodsName,
		// "start_time":       params.StartTime.Format("20060102150405"),
		// "end_time":         params.EndTime.Format("20060102150405"),
		"service_location": params.ServiceLocation,
		"deposit_amount":   params.DepositAmount,
		"rent_unit_fee":    params.RentUnitFee,
		"rent_fee_desc":    params.RentFeeDesc,
	}

	for _, it := range options {
		if it == nil {
			continue
		}
		it(p)
	}
	executor := c.BuildExecutor("https://api.mch.weixin.qq.com/wxv/createrentbill", p).UseXML(true).UseTwowayAuthentication(true)
	// TODO: 智能创建
	results = &PayScoreRentCreateResults{
		Results:    executor.Execute(false),
		Params:     p,
		OutOrderNo: params.OutOrderNo,
		ServiceID:  c.ServiceID,
		MchID:      c.MchID,
	}
	err = results.Error()
	backoff := &Backoff{Max: time.Second * 5}
	for err == nil && results.SubCode() == "SYSTEMERROR" && backoff.Wait() {
		results.Results = executor.Execute(false)
		err = results.Error()
	}
	if err != nil {
		return
	}
	GetLoggerFromContext(ctx).Info(results)
	if !results.Success() {
		code := results.SubCode()
		if code == "" {
			code = "gateway_error"
		}
		err = utils.Err(code, results.SubMessage())
		return
	}
	results.OrderID = results.Get("order_id").String()
	results.MiniprogramAppid = results.Get("miniprogram_appid").String()
	results.MiniprogramPath = results.Get("miniprogram_path").String()
	results.MiniprogramUsername = results.Get("miniprogram_username").String()
	results.Package = results.Get("package").String()
	r, err := c.PayScoreRentQuery(ctx, results.OutOrderNo)
	if err != nil {
		return results, err
	}
	results.State = r.State
	results.FinishTicket = r.FinishTicket
	results.FinishTransactionID = r.FinishTransactionID
	return
}

type Backoff struct {
	Max   time.Duration
	Start time.Duration
	crt   int64
}

func (b *Backoff) Wait() bool {
	if b.Start <= 0 {
		b.Start = 1
	}
	if b.crt <= 0 {
		b.crt = 1
	}
	crt := time.Millisecond * time.Duration(b.crt)
	if b.Max <= crt {
		return false
	}
	//fmt.Println(crt)
	time.Sleep(crt)
	b.crt <<= 1
	return true
}

// PayScoreRentQuery 查询租借订单。接口文档：https://pay.weixin.qq.com/wiki/doc/apiv3/payscore.php?chapter=18_2&index=3
func (c *WechatPayClient) PayScoreRentQuery(ctx context.Context, outOrderNo string, options ...ParamsOption) (results *PayScoreRentCreateResults, err error) {
	// 签名方式：HMAC-SHA256
	p := opensdk.Params{
		"version":      "1.0",
		"sign_type":    "HMAC-SHA256",
		"nonce_str":    opensdk.RandomString(10),
		"mch_id":       c.MchID,
		"appid":        c.AppID,
		"service_id":   c.ServiceID,
		"out_order_no": outOrderNo,
	}

	for _, it := range options {
		if it == nil {
			continue
		}
		it(p)
	}
	executor := c.BuildExecutor("https://api.mch.weixin.qq.com/wxv/queryrentbill", p).UseXML(true).UseTwowayAuthentication(true)
	// TODO: 智能创建
	results = &PayScoreRentCreateResults{
		Results:    executor.Execute(false),
		Params:     p,
		OutOrderNo: outOrderNo,
		ServiceID:  c.ServiceID,
		MchID:      c.MchID,
	}

	err = results.Error()
	// fmt.Println("====ssssssss2==", err)
	// fmt.Println("====ssssssss3==", results.Error())
	backoff := &Backoff{Max: time.Second * 5}
	for err == nil && results.SubCode() == "SYSTEMERROR" && backoff.Wait() {
		results.Results = executor.Execute(false)
		// GetLoggerFromContext(ctx).Error(results.Results, string(results.Results.Body()))
		// fmt.Println("====ssssssss==", results.Results)
		// fmt.Println("====ssssssss==", string(results.Results.Body()))
		err = results.Error()
	}
	if err != nil {
		return
	}
	GetLoggerFromContext(ctx).Info(results.Results, string(results.Results.Body()))
	if !results.Success() {
		code := results.SubCode()
		if code == "" {
			code = "gateway_error"
		}

		err = utils.Err(code, results.SubMessage())
		return
	}
	results.State = results.Get("state").String()
	results.FinishTicket = results.Get("finish_ticket").String()
	results.FinishTransactionID = results.Get("finish_transaction_id").String()
	results.OrderID = results.Get("order_id").String()

	return
}

// PayScoreRentCancel 撤销租借订单。接口文档：https://pay.weixin.qq.com/wiki/doc/apiv3/payscore.php?chapter=18_3&index=4
func (c *WechatPayClient) PayScoreRentCancel(ctx context.Context, outOrderNo, reason string, options ...ParamsOption) (results opensdk.Results, err error) {
	p := opensdk.Params{
		"version":      "1.0",
		"sign_type":    "HMAC-SHA256",
		"nonce_str":    opensdk.RandomString(10),
		"mch_id":       c.MchID,
		"appid":        c.AppID,
		"service_id":   c.ServiceID,
		"out_order_no": outOrderNo,
		"reason":       reason,
	}

	for _, it := range options {
		if it == nil {
			continue
		}
		it(p)
	}
	executor := c.BuildExecutor("https://api.mch.weixin.qq.com/wxv/cancelbill", p).UseXML(true).UseTwowayAuthentication(true)

	results = executor.Execute(false)
	err = results.Error()
	backoff := &Backoff{Max: time.Second * 5}
	for err == nil && results.SubCode() == "SYSTEMERROR" && backoff.Wait() {
		results = executor.Execute(false)
		err = results.Error()
	}
	if err != nil {
		return
	}
	GetLoggerFromContext(ctx).Info(results)
	if !results.Success() && results.SubCode() != "ORDERNOTEXIST" { // 订单未找到也默认成功
		code := results.SubCode()
		if code == "" {
			code = "gateway_error"
		}
		msg := results.SubMessage()
		if msg == "" {
			msg = results.Message()
		}
		//撤销失败，请确认入参
		if code == "PARAM_ERROR" && msg == "撤销失败，请确认入参" {
			GetLoggerFromContext(ctx).Error(utils.Err(code, msg+":", outOrderNo))
			return
		}
		err = utils.Err(code, msg)
		return
	}

	return
}

// 请求参数，必填
type PayScoreRentFinishParams struct {
	OutOrderNo         string    `json:"out_order_no" xml:"out_order_no"`                 // 商户服务订单号
	Returned           bool      `json:"returned" xml:"returned"`                         // 是否归还
	FinishTicket       string    `json:"finish_ticket,omitempty"`                         // 完结凭证
	RealEndTime        time.Time `json:"real_end_time" xml:"real_end_time"`               // 实际归还时间(未归还不填)
	ServiceEndLocation string    `json:"service_end_location" xml:"service_end_location"` // 归还地点(未归还不填)
	TotalAmount        int64     `json:"total_amount" xml:"total_amount"`                 // 总额
	// RentFee             int64     `json:"rent_fee" xml:"rent_fee"`                           // 租金费用
	// CompensationFee     int64     `json:"compensation_fee" xml:"compensation_fee"`           // 赔偿金费用(与 租金费用 二选一)
	CompensationFeeDesc string `json:"compensation_fee_desc" xml:"compensation_fee_desc"` // 赔偿金费用说明
}

// PayScoreRentFinish 完结租借订单。接口文档：https://pay.weixin.qq.com/wiki/doc/apiv3/payscore.php?chapter=18_4&index=5
func (c *WechatPayClient) PayScoreRentFinish(ctx context.Context, params PayScoreRentFinishParams, options ...ParamsOption) (results *PayScoreRentCreateResults, err error) {
	// 签名方式：HMAC-SHA256
	p := opensdk.Params{
		"version":       "1.0",
		"sign_type":     "HMAC-SHA256",
		"nonce_str":     opensdk.RandomString(10),
		"mch_id":        c.MchID,
		"appid":         c.AppID,
		"service_id":    c.ServiceID,
		"out_order_no":  params.OutOrderNo,
		"total_amount":  params.TotalAmount,
		"finish_ticket": params.FinishTicket,
	}
	if params.Returned {
		p["returned"] = "TRUE"
		p["real_end_time"] = params.RealEndTime.Format("20060102150405")
		p["service_end_location"] = params.ServiceEndLocation
		p["rent_fee"] = params.TotalAmount
	} else {
		p["returned"] = "FALSE"
		p["compensation_fee"] = params.TotalAmount
		p["compensation_fee_desc"] = params.CompensationFeeDesc
		p["rent_fee"] = 0
	}
	// if params.CompensationFee > 0 {
	// 	p["compensation_fee"] = params.CompensationFee
	// 	p["compensation_fee_desc"] = params.CompensationFeeDesc
	// }
	// if params.RentFee > 0 {
	// 	p["rent_fee"] = params.RentFee
	// }

	for _, it := range options {
		if it == nil {
			continue
		}
		it(p)
	}
	executor := c.BuildExecutor("https://api.mch.weixin.qq.com/wxv/finishrentbill", p).UseXML(true).UseTwowayAuthentication(true)
	// TODO: 智能创建
	results = &PayScoreRentCreateResults{
		Results:    executor.Execute(false),
		Params:     p,
		OutOrderNo: params.OutOrderNo,
		ServiceID:  c.ServiceID,
		MchID:      c.MchID,
	}
	err = results.Error()
	backoff := &Backoff{Max: time.Second * 5}
	for err == nil && results.SubCode() == "SYSTEMERROR" && backoff.Wait() {
		results.Results = executor.Execute(false)
		err = results.Error()
	}
	if err != nil {
		return
	}
	GetLoggerFromContext(ctx).Info(results.Results, string(results.Results.Body()))

	if !results.Success() {
		code := results.SubCode()
		if code == "" {
			code = "gateway_error"
		}
		err = utils.Err(code, results.SubMessage())
		return
	}

	return
}

func GetLoggerFromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value("logger").(Logger); ok && l != nil {
		return l
	}
	return dLogger
}

type Logger interface {
	Info(a ...interface{})
	Error(a ...interface{})
}

var dLogger = &dropLogger{}

type dropLogger struct {
}

func (*dropLogger) Info(a ...interface{}) {
	logs.Debug(a...)
}
func (*dropLogger) Error(a ...interface{}) {
	logs.Error(a...)
}
