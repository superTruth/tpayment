package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/basekey"
	"tpayment/models"
	"tpayment/models/payment/record"

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
			CardExpDay:     "25",
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
			CardExpDay:     "25",
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

	fmt.Println("rep->", string(repByte))
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

func TestDB(t *testing.T) {
	conf.InitConfigData()

	models.InitDB()

	basekey.Init() // 初始化基础秘钥

	recordBean := &record.TxnRecord{
		BaseModel: models.BaseModel{
			Db: models.DB(),
		},
	}

	total, err := recordBean.GetSettlementTotal(9, 0, 1)

	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Println("total len->", len(total))

	for i, totalTmp := range total {
		fmt.Println(i, "->", *totalTmp)
	}

}
