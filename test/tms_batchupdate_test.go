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

func TestCreateBatchUpdate(t *testing.T) {
	TestLogin(t)

	fmt.Println("add user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	tags := []*tms.DeviceTag{
		{
			BaseModel: models.BaseModel{
				ID: 1,
			},
		},
		{
			BaseModel: models.BaseModel{
				ID: 2,
			},
		},
	}

	deviceModels := []*tms.DeviceModel{
		{
			BaseModel: models.BaseModel{
				ID: 1,
			},
		},
		{
			BaseModel: models.BaseModel{
				ID: 2,
			},
		},
	}

	reqBean := &tms.BatchUpdate{
		Description:  "Test1",
		ConfigTags:   tags,
		ConfigModels: deviceModels,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsBatchUpdateAdd, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateBatchUpdate(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	tags := []*tms.DeviceTag{
		{
			BaseModel: models.BaseModel{
				ID: 1,
			},
		},
	}

	deviceModels := []*tms.DeviceModel{
		{
			BaseModel: models.BaseModel{
				ID: 3,
			},
		},
	}

	reqBean := &tms.BatchUpdate{
		BaseModel: models.BaseModel{
			ID: 3,
		},
		Description:  "Test2",
		ConfigTags:   tags,
		ConfigModels: deviceModels,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsBatchUpdateUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryBatchUpdate(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset: 0,
		Limit:  100,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsBatchUpdateQuery, time.Second*10)

	formatJson(repByte)
}

func TestDeleteBatchUpdate(t *testing.T) {
	TestLogin(t)

	fmt.Println("delete user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 2}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsBatchUpdateDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestStartHandleBatchUpdate(t *testing.T) {
	TestLogin(t)

	fmt.Println("handle user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 23}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsBatchUpdateStartHandle, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
