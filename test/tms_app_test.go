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

func TestCreateTmsApp(t *testing.T) {
	TestLogin(t)

	fmt.Println("add user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.App{
		Name:        "Fang Apk",
		PackageId:   "com.bindo.test",
		Description: "测试",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppAdd, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateTmsApp(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.AppInDevice{
		BaseModel: models.BaseModel{
			ID: 1,
		},
		Name: "Test23423",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInDeviceUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryTmsApp(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset:   0,
		Limit:    100,
		DeviceId: 1505156075807081495,
		//Filters: map[string]string{
		//	"device_sn": "PAX-A920-0821157228",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInDeviceQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteTmsApp(t *testing.T) {
	TestLogin(t)

	fmt.Println("delete user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 1505221047908102133}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInDeviceDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
