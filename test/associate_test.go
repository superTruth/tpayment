package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/models/merchant"
	"tpayment/modules"
)

func TestAddMerchantAssociate(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddMerchantAssociate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.UserMerchantAssociate{
		MerchantId: 4,
		UserId:     5,
		Role:       string(conf.RoleAdmin),
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteMerchantAssociate(t *testing.T) {
	TestLogin(t)

	fmt.Println("DeleteMerchantAssociate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID:1}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryMerchantAssociate(t *testing.T) {
	token := Login("fang.qiang3@bindo.com", "123456")

	fmt.Println("query user", line)
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

	repByte,_ := post(reqByte, header, "http://localhost:80/payment/account/query", time.Second*10)

	fmt.Println("rep->", string(repByte))

}

func TestQueryUserInMerchant(t *testing.T) {
	token := Login("fang.qiang3@bindo.com", "123456")
	fmt.Println("query user in merchant", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		MerchantId: 3,
		Offset:  0,
		Limit:   100,
		//Filters: map[string]string{
		//	"pwd": "123456",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte,_ := post(reqByte, header, BaseUrl+conf.UrlQueryUserInMerchantQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
