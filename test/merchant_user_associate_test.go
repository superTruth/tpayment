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
		UserId:     8,
		Role:       string(conf.RoleAdmin),
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteMerchantAssociate(t *testing.T) {
	TestLogin(t)

	fmt.Println("DeleteMerchantAssociate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 4}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryMerchantAssociate(t *testing.T) {
	token := Login("fang.qiang6@bindo.com", "123456")

	fmt.Println("query user", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		MerchantId: 4,
		Offset:     0,
		Limit:      100,
		//Filters: map[string]string{
		//	"pwd": "123456",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))

}
