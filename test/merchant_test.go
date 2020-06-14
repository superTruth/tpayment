package test

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
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
		AgencyId: 4,
		Name: "merchant 1",
		Tel:  "123456789",
		Addr: "wuxicun",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlMerchantAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateMerchant(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddMerchant", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.Merchant{
		Model: gorm.Model{
			ID:        3,
		},
		Name: "merc",
		Tel:  "",
		Addr: "wuxicun2",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlMerchantUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryMerchant(t *testing.T) {
	TestLogin(t)
	fmt.Println("query merchant", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		AgencyId: 4,
		Offset:  0,
		Limit:   100,
		//Filters: map[string]string{
		//	"pwd": "123456",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlMerchantQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
