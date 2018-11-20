package alipay

import (
	"encoding/json"
	"fmt"
	"time"
)

type JsonTime time.Time

//MarshalJSON 实现它的json序列化方法
func (t JsonTime) MarshalJSON() ([]byte, error) {
	s := time.Time(t).Format("2006-01-02 15:04:05")
	return []byte("\"" + s + "\""), nil
}

//MarshalJSON ,UnmarshalJSON , String

// MiniProgramCreditBorrowParams 小程序芝麻信用借还参数。该参数适用于小程序内下单。
// 参数字符串的某一项值包含中文，请使用encodeURIComponent对该值进行编码
type MiniProgramCreditBorrowParams struct {
	OutOrderNo      string   `json:"out_order_no,omitempty"`
	ProductCode     string   `json:"product_code,omitempty"`
	GoodsGame       string   `json:"goods_name,omitempty"`
	RentUnit        string   `json:"rent_unit,omitempty"`
	RentAmount      float64  `json:"rent_amount,omitempty"`
	DepositAmount   float64  `json:"deposit_amount,omitempty"`
	DepositState    string   `json:"deposit_state,omitempty"`
	InvokeReturnURL string   `json:"invoke_return_url,omitempty"`
	InvokeType      string   `json:"invoke_type,omitempty"`
	CreditBiz       string   `json:"credit_biz,omitempty"`
	BorrowTime      JsonTime `json:"borrow_time,omitempty"`
	ExpiryTime      JsonTime `json:"expiry_time,omitempty"`
	MobileNo        string   `json:"mobile_no,omitempty"`
	BorrowShopName  string   `json:"borrow_shop_name,omitempty"`
	RentSettleType  string   `json:"rent_settle_type,omitempty"`
	InvokeState     string   `json:"invoke_state,omitempty"`
	RentInfo        string   `json:"rent_info,omitempty"`
	Name            string   `json:"name,omitempty"`
	CertNo          string   `json:"cert_no,omitempty"`
	Address         string   `json:"address,omitempty"`
}

func (p *MiniProgramCreditBorrowParams) String() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//BuildMiniProgramCreditBorrowParams 构建信用借还的参数。接口文档：https://docs.alipay.com/mini/api/zmcreditborrow
func (c *Client) BuildMiniProgramCreditBorrowParams(outOrderNo, goodsName string,
	rentUnit string, rentAmount, depositAmount float64,
	invokeReturnPage, creditBiz string,
	borrowTime, expiryTime time.Time) *MiniProgramCreditBorrowParams {

	if !(rentUnit == "DAY_YUAN" || rentUnit == "HOUR_YUAN" || rentUnit == "YUAN" || rentUnit == "YUAN_ONCE") {
		panic("rentUnit 参数值错误")
	}

	return &MiniProgramCreditBorrowParams{
		OutOrderNo:      outOrderNo,
		ProductCode:     "w1010100000000002858",
		GoodsGame:       goodsName,
		RentUnit:        rentUnit,
		RentAmount:      rentAmount,
		DepositAmount:   depositAmount,
		DepositState:    "Y",
		InvokeReturnURL: fmt.Sprintf("alipays://platformapi/startapp?appId=%s&page=%s", c.AppID, invokeReturnPage),
		InvokeType:      "TINYAPP",
		CreditBiz:       creditBiz,
		BorrowTime:      JsonTime(borrowTime),
		ExpiryTime:      JsonTime(expiryTime),
	}
}

type MerchantOrderRentQueryResponse struct {
	OrderNo           string `json:"order_no,omitempty"`
	GoodsName         string `json:"goods_name,omitempty"`
	UserID            string `json:"user_id,omitempty"`
	BorrowTime        string `json:"borrow_time,omitempty"`
	RestoreTime       string `json:"restore_time,omitempty"`
	UseState          string `json:"use_state,omitempty"`
	PayStatus         string `json:"pay_status,omitempty"`
	PayAmountType     string `json:"pay_amount_type,omitempty"`
	PayAmount         string `json:"pay_amount,omitempty"`
	PayTime           string `json:"pay_time,omitempty"`
	AdmitState        string `json:"admit_state,omitempty"`
	AlipayFundOrderNo string `json:"alipay_fund_order_no,omitempty"`
	Code              string `json:"code,omitempty"`
}

// MerchantOrderRentQuery 芝麻信用借还提供的供商户查询借还订单详情。接口文档： https://docs.open.alipay.com/api_8/zhima.merchant.order.rent.query
func (c *Client) MerchantOrderRentQuery(outOrderNo string) *Executor {

	params := BuildParams("zhima.merchant.order.rent.query")

	// biz := map[string]string{
	// 	"out_order_no": outOrderNo,
	// 	"product_code": "w1010100000000002858",
	// }
	params["out_order_no"] = outOrderNo
	params["product_code"] = "w1010100000000002858"

	return c.Execute(params, "zhima_merchant_order_rent_query_response")

}

/*
state
:
"Y"
invoke_state
:
undefined
order_no
:
"10020027981"
out_order_no
:
"2018101761217333128"
resultStatus
:
9000
user_id
:
"2088002463165793"
*/
