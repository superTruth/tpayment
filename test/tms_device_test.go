package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/models/tms"
	"tpayment/modules"

	"github.com/jinzhu/gorm"
)

func TestUpdateTmsDevice(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	var deviceTags []tms.DeviceTagFull
	_ = json.Unmarshal([]byte(`[
					{
						"ID": 25,
						"Name": "tag1"
					},
					{
						"AgencyId": "0",
						"CreatedAt": "0001-01-01T00:00:00Z",
						"DeletedAt": null,
						"ID": 24,
						"MidId": 744,
						"Name": "tag2",
						"UpdatedAt": "0001-01-01T00:00:00Z"
					}
				]`), &deviceTags)

	reqBean := &tms.DeviceInfo{
		Model: gorm.Model{
			ID: 1505156075807081492,
		},
		DeviceCsn: "12312",
		Tags:      deviceTags,
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
		Filters: map[string]string{
			"device_sn": "PAX-A920-0821157228",
		},
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

	reqBean := &modules.BaseIDRequest{ID: 1505156075807081492}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsDeviceDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
