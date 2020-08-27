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

func TestCreateTmsAppInBatchUpdate(t *testing.T) {
	TestLogin(t)

	fmt.Println("add user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.AppInDevice{
		ExternalId: 7,
		Status:     conf.TmsStatusPendingInstall,
		AppID:      7,
		AppFileId:  12,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInBatchUpdateAdd, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateTmsAppInBatchUpdate(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.AppInDevice{
		BaseModel: models.BaseModel{
			ID: 44,
		},
		Status: conf.TmsStatusPendingUninstalled,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInBatchUpdateUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryTmsAppInBatchUpdate(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset:  0,
		Limit:   100,
		BatchId: 3,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInBatchUpdateQuery, time.Second*10)

	fmt.Println("rep->", string(repByte))

	formatJson(repByte)
}

func TestDeleteTmsAppInBatchUpdate(t *testing.T) {
	TestLogin(t)

	fmt.Println("delete user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 44}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppInBatchUpdateDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
