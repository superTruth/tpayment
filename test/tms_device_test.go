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

func TestUpdateTmsDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	var deviceTags []*tms.DeviceTagFull
	_ = json.Unmarshal([]byte(`[
					{
						"id": 25
					}
				]`), &deviceTags)

	fmt.Println("deviceTags len->", len(deviceTags))

	reqBean := &tms.DeviceInfo{
		BaseModel: models.BaseModel{
			ID: 1,
		},
		DeviceCsn: "456789",
		Tags:      &deviceTags,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsDeviceUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryTmsDevice(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset: 0,
		Limit:  100,
		//Filters: map[string]string{
		//	"device_sn": "PAX-A920-0821157228",
		//},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsDeviceQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))
	formatJson(repByte)
	//fmt.Println("rep->", string(repByte))
}

func TestDeleteTmsDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("delete user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 1}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsDeviceDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
