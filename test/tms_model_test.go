package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
)

func TestAddTmsModel(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestAddTmsTag", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.DeviceTag{
		Name: "Tag4",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsModelAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateTmsModel(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.DeviceTag{
		BaseModel: models.BaseModel{ID: 27},
		Name:      "Tag444",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsModelUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryTmsModel(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset: 0,
		Limit:  100,
		Filters: map[string]string{
			"name": "tag",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsModelQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))
	formatJson(repByte)
	//fmt.Println("rep->", string(repByte))
}

func TestDeleteTmsModel(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestDeleteTmsTag", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 26}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsModelDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
