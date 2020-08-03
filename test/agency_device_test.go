package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/modules"
	"tpayment/modules/agency/agencydevice"
)

func TestAddDeviceAcquirerByFile(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestAddDeviceAcquirer", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &agencydevice.DeviceBindRequest{
		AgencyId: 4,
		DeviceId: 0,
		FileUrl:  "https://mdmfiles.oss-cn-hongkong.aliyuncs.com/other%20file/test.csv",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyDeviceAdd, time.Second*10)

	formatJson(repByte)
}

func TestAddDeviceAcquirerByID(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestAddDeviceAcquirer", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &agencydevice.DeviceBindRequest{
		AgencyId: 4,
		DeviceId: 3,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyDeviceAdd, time.Second*10)

	formatJson(repByte)
}

func TestQueryDeviceAcquirer(t *testing.T) {
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
			"device_sn": "PAX-",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyDeviceQuery, time.Second*10)

	formatJson(repByte)
}

func TestDeleteDeviceAcquirer(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 3}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAgencyDeviceDelete, time.Second*10)

	formatJson(repByte)
}
