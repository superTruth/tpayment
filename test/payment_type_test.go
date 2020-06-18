package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/modules"
)

func TestQueryPaymentTypes(t *testing.T) {
	TestLogin(t)
	fmt.Println("query payment types", line)
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

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlAgencyPaymentTypes, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryPaymentMethods(t *testing.T) {
	TestLogin(t)
	fmt.Println("query payment methods", line)
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

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlAgencyPaymentMethods, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryEntryTypes(t *testing.T) {
	TestLogin(t)
	fmt.Println("query payment methods", line)
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

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlAgencyEntryTypes, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
