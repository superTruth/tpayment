package test

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/modules"
)

func TestAddAgency(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddAgency", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &agency.Agency{
		Name: "agency 1",
		Tel:  "123456789",
		Addr: "wuxicun",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlAgencyAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateAgency(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateAgency", line)

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

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlAgencyUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryAgency(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset:  0,
		Limit:   100,
		//Filters: map[string]string{
		//	"pwd": "123456",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlAgencyQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
