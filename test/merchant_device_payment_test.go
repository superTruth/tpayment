package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
)

func TestAddMerchantPaymentDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestAddMerchantPaymentDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.PaymentSettingInDevice{
		MerchantDeviceId: 2,
		PaymentMethods: &models.StringArray{
			"Visa", "MasterCard", "Unionpay",
		},
		EntryTypes: &models.StringArray{
			"Swipe", "Contact", "Contactless",
		},
		PaymentTypes: &models.StringArray{
			"Sale", "Void", "Refund",
		},
		AcquirerId: 1,
		Mid:        "123456789012345",
		Tid:        "12345678",
		Addition:   "http://test.com",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDevicePaymentAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateMerchantPaymentDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddMerchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.PaymentSettingInDevice{
		BaseModel: models.BaseModel{
			ID: 1,
		},
		PaymentMethods: &models.StringArray{},
		EntryTypes: &models.StringArray{
			"Swipe", "Contact",
		},
		PaymentTypes: &models.StringArray{
			"Sale", "Refund",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDevicePaymentUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteMerchantPaymentDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("Delete Merchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.PaymentSettingInDevice{
		BaseModel: models.BaseModel{
			ID: 1,
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDevicePaymentDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryMerchantPaymentDevice(t *testing.T) {
	TestLogin(t)
	fmt.Println("query merchant", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		DeviceId: 2,
		Offset:   0,
		Limit:    100,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDevicePaymentQuery, time.Second*10)

	formatJson(repByte)
}
