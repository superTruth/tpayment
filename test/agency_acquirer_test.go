package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/modules"
)

func TestAddAgencyAcquirer(t *testing.T) {
	TestLogin(t)

	fmt.Println("AddAgencyAcquirer", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &agency.Acquirer{
		Name:          "BOC",
		Addition:      "addtion",
		ConfigFileUrl: "https://asdfadf",
		AgencyId:      4,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyAcquirerAdd, time.Second*10)

	formatJson(repByte)
}

func TestUpdateAgencyAcquirer(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateAgencyAcquirer", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &agency.Acquirer{
		BaseModel: models.BaseModel{
			ID: 3,
		},
		Name:          "BOC3",
		Addition:      "addtion3",
		ConfigFileUrl: "https://asdfadf3",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyAcquirerUpdate, time.Second*10)

	formatJson(repByte)
}

func TestQueryAgencyAcquirer(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset:   0,
		Limit:    100,
		AgencyId: 4,
		Filters: map[string]string{
			"name": "BOC",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyAcquirerQuery, time.Second*10)

	formatJson(repByte)
}

func TestDeleteAgencyAcquirer(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &agency.Acquirer{
		BaseModel: models.BaseModel{
			ID: 3,
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyAcquirerDelete, time.Second*10)

	formatJson(repByte)
}
