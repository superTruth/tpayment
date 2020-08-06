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

func TestAddTmsTag(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestAddTmsTag", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.DeviceTag{
		Name: "Tag4",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsTagAdd, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateTmsTag(t *testing.T) {
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

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsTagUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryTmsTag(t *testing.T) {
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

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsTagQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))
	formatJson(repByte)
	//fmt.Println("rep->", string(repByte))
}

func TestDeleteTmsTag(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestDeleteTmsTag", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 26}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsTagDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
