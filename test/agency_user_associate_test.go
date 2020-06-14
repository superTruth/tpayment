package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/modules"
)

func TestAddAgencyAssociate(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddAgencyAssociate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &agency.UserAgencyAssociate{
		AgencyId: 4,
		UserId:   7,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyAssociateAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteAgencyAssociate(t *testing.T) {
	TestLogin(t)

	fmt.Println("DeleteAgencyAssociate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 1}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyAssociateDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryAgencyAssociate(t *testing.T) {
	token := Login("fang.qiang6@bindo.com", "123456")

	fmt.Println("query user", line)
	fmt.Println("token->", token)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		AgencyId: 4,
		Offset:   0,
		Limit:    100,
		//Filters: map[string]string{
		//	"pwd": "123456",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyAssociateQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))

}
