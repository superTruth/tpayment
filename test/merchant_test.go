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

func TestAddMerchant(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddMerchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.Merchant{
		AgencyId: 7,
		Name:     "merchant 1",
		Tel:      "123456789",
		Addr:     "wuxicun",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateMerchant(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddMerchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.Merchant{
		BaseModel: models.BaseModel{
			ID: 9,
		},
		Name: "merc",
		Tel:  "",
		Addr: "wuxicun2",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteMerchant(t *testing.T) {
	TestLogin(t)

	fmt.Println("Delete Merchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.Merchant{
		BaseModel: models.BaseModel{
			ID: 9,
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryMerchant(t *testing.T) {
	TestLogin(t)
	fmt.Println("query merchant", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset: 0,
		Limit:  100,
		Filters: map[string]string{
			"name": "merchant",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantQuery, time.Second*10)

	formatJson(repByte)
}

func TestQueryMerchantInAgency(t *testing.T) {
	TestLogin(t)
	fmt.Println("query merchant in agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		AgencyId: 4,
		Offset:   0,
		Limit:    100,
		Filters: map[string]string{
			"name": "merchant",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantQuery, time.Second*10)

	formatJson(repByte)
}
