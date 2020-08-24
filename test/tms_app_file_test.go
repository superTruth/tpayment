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

func TestCreateTmsAppFile(t *testing.T) {
	TestLogin(t)

	fmt.Println("add user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	// https://mdmfiles.oss-cn-hongkong.aliyuncs.com/other%20file/Pax-MDM-V1.14_alpha_release_sign.apk
	reqBean := &tms.AppFile{
		UpdateDescription: "MDM First Time",
		FileUrl:           "https://mdmfiles.oss-cn-hongkong.aliyuncs.com/other%20file/Landi-MDM-V1.15_alpha_release_20200720%20%281%29.apk",
		AppId:             1,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppFileAdd, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func TestUpdateTmsAppFile(t *testing.T) {
	TestLogin(t)

	fmt.Println("TestUpdateDevice", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &tms.App{
		BaseModel: models.BaseModel{
			ID: 1,
		},
		Name:        "MDM",
		PackageId:   "com.bindo.mdm",
		Description: "MDM",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppUpdate, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryTmsAppFile(t *testing.T) {
	TestLogin(t)
	fmt.Println("query agency", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		AppId:  5,
		Offset: 0,
		Limit:  100,
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppFileQuery, time.Second*10)

	formatJson(repByte)
}

func TestDeleteTmsAppFile(t *testing.T) {
	TestLogin(t)

	fmt.Println("delete user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 1}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlTmsAppDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}
