package test

import (
	"encoding/json"
	"fmt"
	"github.com/go-gomail/gomail"
	"testing"
)

type DefType struct {
	StrArray []string `json:"strs"`
}

func TestMyFunc(t *testing.T) {
	demoData := `{"strs":["12345", "56789"]}`

	demoStruct := &DefType{}

	err := json.Unmarshal([]byte(demoData), demoStruct)
	if err != nil {
		t.Error("err->", err.Error())
		return
	}

	fmt.Println("json->", demoStruct.StrArray)
}

func TestEmail(t *testing.T) {
	m := gomail.NewMessage()

	m.SetHeader("To", "446876407@qq.com")
	m.SetAddressHeader("From", "fang.qiang@bindo.com", "no_reply")
	m.SetHeader("Subject", "No Theme")

	body := `Active here <a href = "https://www.latelee.org">Click</a><br>`
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.office365.com", 587, "fang.qiang@bindo.com", "F1a2n3g4")
	err := d.DialAndSend(m)

	if err != nil {
		t.Error("send error->", err.Error())
		return
	}
}
