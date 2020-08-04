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

func TestCreateTmsAppInDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("add user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.AppInDevice{
		ExternalId:  1,
		Name:        "Test",
		PackageId:   "com.truth.test",
		VersionName: "v1.0.1",
		VersionCode: 123,
		Status:      conf.TmsStatusWarningInstalled,
		AppID:       0,
		AppFileId:   0,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInDeviceAdd, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateTmsAppInDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.AppInDevice{
		BaseModel: models.BaseModel{
			ID: 1505221047908102138,
		},
		Status: conf.TmsStatusPendingUninstalled,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInDeviceUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryTmsAppInDevice(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset:   0,
		Limit:    100,
		DeviceId: 1,
		//Filters: map[string]string{
		//	"device_sn": "PAX-A920-0821157228",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInDeviceQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))

	formatJson(repByte)
}

func TestDeleteTmsAppInDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("delete user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 1505221047908102137}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInDeviceDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
