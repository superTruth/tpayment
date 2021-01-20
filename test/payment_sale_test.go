package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/modules"
	"tpayment/modules/payment/pay_manage"

	"github.com/google/uuid"
)

func TestSaleVisa(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestSaleVisa", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &api_define.TxnReq{
		Uuid:          uuid.New().String(),
		TxnType:       conf.Sale,
		DeviceID:      "A920123456789",
		PaymentMethod: conf.RequestCreditCard,
		MerchantID:    9,
		Amount:        "0.2",
		Currency:      "USD",
		CreditCardBean: &api_define.CreditCardBean{
			CardReaderMode: conf.Contact,
			CardExpMonth:   "12",
			CardExpYear:    "24",
			CardFallback:   false,
			CardNumber:     "4384375620640049",
			CardSn:         "1",
			CardTrack2:     "4384375620640049D24122012000001000872",
			CardHolderName: "test",
			IccRequest:     "5A0843843756206400495F2A0203445F34010182023C008407A0000000031010950500000080009A031811209B02E8009C01009F02060000000001009F03060000000000009F080200969F090200969F1A0203449F1E0831323334353637389F2608E92EDAAE12C246789F2701809F3303E0B0C89F34031E03009F3501229F360209049F37048B4CF2449F4104000000039F100706010A03A02002",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlSale, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestSaleCUP(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestSaleCUP", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &api_define.TxnReq{
		Uuid:          uuid.New().String(),
		TxnType:       conf.Sale,
		DeviceID:      "A920123456789",
		PaymentMethod: conf.RequestCreditCard,
		MerchantID:    9,
		Amount:        "0.01",
		Currency:      "HKD",
		CreditCardBean: &api_define.CreditCardBean{
			CardReaderMode: conf.ContactLess,
			CardExpMonth:   "10",
			CardExpYear:    "30",
			CardFallback:   false,
			CardNumber:     "6210946888140008",
			CardSn:         "01",
			CardTrack2:     "6210946888140008D30102010000000000000",
			CardHolderName: "test",
			PIN:            "111111",
			//IccRequest:     "9F360208EB9F3501229F4104000000019C01009A031912109F3303E068C89F1E0831313730303938329F2701809F101307010103A00000010A0100000000005E4EAD3E9F26084A8025B8F97EBB4E9F02060000000005005F2A0201568408A0000003330101029F370437B9C5A382027C00950500000000009F1A020156",
			IccRequest: "9F26084A8025B8F97EBB4E9F02060000000005009F4104000000019A031912109F3303E068C89F1A0201569F101307010103A00000010A0100000000005E4EAD3E9F1E0831313730303938329C01005F2A0201569F370437B9C5A382027C009F2701809F360208EB8408A0000003330101029F35012295050000000000",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlSale, time.Second*10)

	formatJson(repByte)
}

func TestSaleVisaOffline(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestSaleVisaOffline", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	nowTime := time.Now()
	reqBean := &api_define.TxnReq{
		Uuid:          uuid.New().String(),
		TxnType:       conf.Sale,
		DeviceID:      "A92012345678",
		PaymentMethod: conf.RequestCreditCard,
		MerchantID:    9,
		Amount:        "0.2",
		Currency:      "USD",

		DateTime:           &nowTime,
		AcquirerMerchantID: "123456789012345",
		AcquirerTerminalID: "",
		AcquirerRRN:        "123456789",
		AcquirerType:       conf.Visa,
		AcquirerReconID:    "123456",

		CreditCardBean: &api_define.CreditCardBean{
			CardReaderMode: conf.Contact,
			CardExpMonth:   "12",
			CardExpYear:    "24",
			CardFallback:   false,
			CardNumber:     "*********0049",
			CardSn:         "1",
			CardTrack2:     "4384375620640049D24122012000001000872",
			CardHolderName: "test",
			IccRequest:     "5A0843843756206400495F2A0203445F34010182023C008407A0000000031010950500000080009A031811209B02E8009C01009F02060000000001009F03060000000000009F080200969F090200969F1A0203449F1E0831323334353637389F2608E92EDAAE12C246789F2701809F3303E0B0C89F34031E03009F3501229F360209049F37048B4CF2449F4104000000039F100706010A03A02002",

			IccResponse:  "12331234232",
			AuthCode:     "12345",
			TraceNum:     123,
			BatchNum:     4567,
			ResponseCode: "00",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlSaleOffline, time.Second*10)

	formatJson(repByte)
}

func TestSaleOfflineCustomer(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestSaleVisaOffline", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	nowTime := time.Now()
	reqBean := &api_define.TxnReq{
		Uuid:          uuid.New().String(),
		TxnType:       conf.Sale,
		DeviceID:      "A920123456789",
		PaymentMethod: conf.RequestOther,
		MerchantID:    9,
		Amount:        "0.2",
		Currency:      "USD",

		DateTime:              &nowTime,
		AcquirerRRN:           "123456789",
		AcquirerReconID:       "123456",
		CustomerPaymentMethod: "my payment",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlSaleOffline, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestVoidOffline(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestVoidOffline", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &api_define.TxnReq{
		Uuid:       uuid.New().String(),
		TxnType:    conf.Void,
		MerchantID: 9,
		OrgTxnID:   1908925000841166848,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlVoidOffline, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestRefundOffline(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestRefundOffline", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &api_define.TxnReq{
		Uuid:       uuid.New().String(),
		TxnType:    conf.Refund,
		MerchantID: 9,
		OrgTxnID:   1909932860635086848,
		Amount:     "0.1",
		Currency:   "USD",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlRefundOffline, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestCheck(t *testing.T) {
	TestLogin(t)

	fmt.Println("Test Check", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &pay_manage.CheckRequest{
		MerchantId:  9,
		TxnID:       1924522619676131328,
		PartnerUUID: "",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlCheck, time.Second*10)

	formatJson(repByte)
}

func TestPaymentConfig(t *testing.T) {
	TestLogin(t)
	fmt.Println("TestPaymentConfig", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		MerchantId: 9,
		DeviceSN:   "555555",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlPaymentConfig, time.Second*10)

	formatJson(repByte)
}
