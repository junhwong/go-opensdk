package tests

import (
	"testing"
	"time"

	"github.com/junhwong/go-opensdk/alipay"
	"github.com/junhwong/go-opensdk/opensdk"
)

var res_ok = `{"zhima_merchant_order_rent_query_response":{"pay_amount":"6.00","pay_amount_type":"RENT","use_state":"restore","admit_state":"Y","goods_name":"租借订单","code":"10000","pay_time":"2018-10-21 05:24:55","msg":"Success","order_no":"92506636","pay_status":"PAY_SUCCESS","restore_time":"2018-10-21 05:24:54","borrow_time":"2018-10-21 03:08:37","expiry_time":"2018-11-01 03:08:34","user_id":"3453532323424","alipay_fund_order_no":"20181121220014256756756740132"},"sign":"X5nftiDMbQl2zyJ1tZbTlnl51h98oupKux6Ry5TFvEWzt87JN5Q0ug49wZqDT307245e5tZpE8iJnMj7+vh4bJkckzt7r1RlIzn3LtZBpBBROXwdRPsFStsTjO/4juFlDyNRriL2u9rO46h0P5ePAtyZPJ8sQhXjVLPZGAUrtS+RxLRqhDae6XQBK5qoH5v6SZgn18rmZtiZRI/g0oOo3gzogDyNtFZ1F96SBbpwyThsFeQTNTF+eH5Lekafq5MrSMoSxVPT6xy6RpVoQus2343242YQuivMtfK/3K7ODAAvdzJ7yo343243eTl2FwIZ4UCGQFkIKrTwXHS6xE1Yw=="}`

func TestToString(t *testing.T) {
	params := opensdk.Params{}
	params["app_id"] = "testsdfsd345345fdsfds"
	params["format"] = "json"
	params["charset"] = "utf-8"
	params["version"] = "1.0"
	params["sign_type"] = "RSA2"
	params["num"] = 1
	params["timestamp"] = opensdk.JsonTime(time.Now())
	params["biz_content"] = &alipay.MiniProgramCreditBorrowParams{
		GoodsGame:  "测试商品",
		RentAmount: "45.56",
	}
	t.Log(params.Sort().ToURLParams(true))

}

func TestToJSON(t *testing.T) {
	params := opensdk.Params{}
	params["app_id"] = "testsdfsd345345fdsfds"
	params["format"] = "json"
	params["charset"] = "utf-8"
	params["version"] = "1.0"
	params["sign_type"] = "RSA2"
	params["num"] = 1
	params["timestamp"] = opensdk.Time{time.Now(), ""}
	params["biz_content"] = &alipay.MiniProgramCreditBorrowParams{
		GoodsGame:  "测试商品",
		RentAmount: "45.56",
	}
	t.Log(params.Sort().ToJSON())

}

func TestToXML(t *testing.T) {
	params := opensdk.Params{}
	params["app_id"] = "testsdfsd345345fdsfds"
	params["format"] = "json"
	params["charset"] = "utf-8"
	params["version"] = "1.0"
	params["sign_type"] = "RSA2"
	params["num"] = 1
	params["timestamp"] = opensdk.Time{time.Now(), ""}
	params["biz_content"] = &alipay.MiniProgramCreditBorrowParams{
		GoodsGame:  "测试商品",
		RentAmount: "45.56",
	}
	t.Log(params.Sort().ToXML())

}
