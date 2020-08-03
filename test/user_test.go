package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"tpayment/conf"
	"tpayment/modules"
	"tpayment/modules/user"
)

const BaseUrl = "http://localhost:80"

//const BaseUrl = "http://paymentstg.horizonpay.cn:80"

func post(reqBody []byte, header http.Header, destUrl string, timeOut time.Duration) (respBody []byte, err error) {
	req, err := http.NewRequest("POST", destUrl, bytes.NewBuffer(reqBody))
	req.Header = header
	defer req.Body.Close()

	client := &http.Client{Timeout: timeOut}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("post fail->", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ = ioutil.ReadAll(resp.Body)

	return respBody, nil
}

var (
	token string
	line  = "--------------------------------------"
)

func ParseResponse(resp []byte, data interface{}) error {
	baseResponse := new(modules.BaseResponse)

	baseResponse.Data = data

	_ = json.Unmarshal(resp, &baseResponse)

	return nil
}

func TestLogin(t *testing.T) {
	fmt.Println("login", line)
	reqBean := &user.LoginRequest{
		Email:     "fang.qiang@bindo.com",
		Pwd:       "123456",
		AppId:     "123456",
		AppSecret: "123456",
	}

	// agency
	//reqBean := &user.LoginRequest{
	//	Email:     "fang.qiang7@bindo.com",
	//	Pwd:       "123456",
	//	AppId:     "123456",
	//	AppSecret: "123456",
	//}

	reqByte, _ := json.Marshal(reqBean)
	repByte, _ := post(reqByte, nil, BaseUrl+conf.UrlAccountLogin, time.Second*10)

	respBean := &user.LoginResponse{}

	_ = ParseResponse(repByte, respBean)
	//json.Unmarshal(repByte, respBean)

	token = respBean.Token

	fmt.Println("rep->", string(repByte))
}

func Login(account, pwd string) string {
	reqBean := &user.LoginRequest{
		Email:     account,
		Pwd:       pwd,
		AppId:     "123456",
		AppSecret: "123456",
	}

	reqByte, _ := json.Marshal(reqBean)
	repByte, _ := post(reqByte, nil, BaseUrl+conf.UrlAccountLogin, time.Second*10)

	respBean := &user.LoginResponse{}

	_ = ParseResponse(repByte, respBean)

	fmt.Println("rep->", string(repByte))

	return respBean.Token
}

func TestLogout(t *testing.T) {
	TestLogin(t)

	fmt.Println("logout", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	repByte, _ := post(nil, header, BaseUrl+conf.UrlAccountLogout, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func TestValidate(t *testing.T) {
	TestLogin(t)

	fmt.Println("validate", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	repByte, _ := post(nil, header, BaseUrl+conf.UrlAccountValidate, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	formatJson(repByte)
}

func TestAddUser(t *testing.T) {
	TestLogin(t)
	//token := Login("fang.qiang7@bindo.com", "123456")

	fmt.Println("add user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &user.AddUserRequest{
		Email: "fang.qiang8@bindo.com",
		Pwd:   "123456",
		Role:  string(conf.RoleUser),
		Name:  "Fang",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAccountAdd, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func TestDeleteUser(t *testing.T) {
	TestLogin(t)

	fmt.Println("delete user", line)

	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseIDRequest{ID: 2}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAccountDelete, time.Second*10)

	fmt.Println("rep->", string(repByte))
}

func TestQueryUser(t *testing.T) {
	TestLogin(t)

	fmt.Println("query user", line)
	header := http.Header{
		conf.HeaderTagToken: []string{token},
	}

	reqBean := &modules.BaseQueryRequest{
		Offset: 0,
		Limit:  100,
		Filters: map[string]string{
			"email": "fang.qiang",
		},
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, header, BaseUrl+conf.UrlAccountQuery, time.Second*10)

	formatJson(repByte)
}

func TestRegister(t *testing.T) {
	reqBean := &user.AddUserRequest{
		Email: "fang.qiang2@bindo.com",
		Pwd:   "123456",
		Name:  "Fang",
	}

	reqByte, _ := json.Marshal(reqBean)

	repByte, _ := post(reqByte, nil, BaseUrl+conf.UrlAccountRegister, time.Second*10)

	respBean := &modules.BaseResponse{}

	_ = json.Unmarshal(repByte, respBean)

	fmt.Println("rep->", string(repByte))
}

func formatJson(str []byte) {
	var dataMap interface{}
	_ = json.Unmarshal(str, &dataMap)

	ret, _ := json.MarshalIndent(dataMap, "", "\t")

	fmt.Println("formatJson->", string(ret))
}
