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

func TestAddMerchantAssociate(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddMerchantAssociate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.UserMerchantAssociate{
		MerchantId: 6,
		UserId:     78,
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

	reqBean := &modules.BaseIDRequest{ID: 8}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateMerchantAssociate(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateMerchantAssociate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &merchant.UserMerchantAssociate{
		BaseModel: models.BaseModel{
			ID: 8,
		},
		Role: string(conf.RoleUser),
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryMerchantAssociate(t *testing.T) {
	//token := Login("fang.qiang6@bindo.com", "123456")
	TestLogin(t)

	fmt.Println("query user", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		MerchantId: 8,
		Offset:     0,
		Limit:      100,
		//Filters: map[string]string{
		//	"pwd": "123456",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlMerchantAssociateQuery, time.Second*10)

	formatJson(repByte)

}
