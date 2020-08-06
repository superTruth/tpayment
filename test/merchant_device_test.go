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

func TestAddMerchantDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestAddMerchantDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.DeviceInMerchant{
		DeviceId:   2,
		MerchantId: 8,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDeviceAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateMerchantDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddMerchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.DeviceInMerchant{
		BaseModel: models.BaseModel{
			ID: 1,
		},
		DeviceId:   1,
		MerchantId: 987,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDeviceUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteMerchantDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("Delete Merchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.DeviceInMerchant{
		BaseModel: models.BaseModel{
			ID: 1,
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDeviceDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryMerchantDevice(t *testing.T) {
	TestLogin(t)
	fmt.Println("query merchant-", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		MerchantId: 8,
		Offset:     0,
		Limit:      100,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDeviceQuery, time.Second*10)

	formatJson(repByte)
}
