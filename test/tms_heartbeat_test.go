package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/modules/tms/clientapi"
)

func TestHeartBeat(t *testing.T) {
	TestLogin(t)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	apps := []*clientapi.AppInfo{
		{
			Name:        "Test1",
			PackageId:   "com.truth.test1",
			VersionName: "v1.0",
			VersionCode: 1,
		},
		{
			Name:        "Test2",
			PackageId:   "com.truth.test2",
			VersionName: "v1.0",
			VersionCode: 1,
		},
	}
	reqBean := &clientapi.RequestBean{
		DeviceSn:    "98210613995100",
		LocationLat: "123",
		LocationLon: "456",
		DeviceModel: "A920",
		Battery:     99,
		AppInfos:    apps,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsHeartBeat, time.Second*10)

	formatJson(repByte)
}
