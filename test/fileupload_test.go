package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/modules/fileupload"
)

func TestFileUpload(t *testing.T) {
	TestLogin(t)

	fmt.Println("FileUpload", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &fileupload.UploadFileRequest{
		FileName: "test1",
		FileSize: 1000000,
		Md5:      "1B2M2Y8AsgTpgAmY7PhCfg==",
		Tag:      "tms",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlFileAdd, time.Second*10)
	formatJson(repByte)
}
