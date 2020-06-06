package test

import (
	"encoding/json"
	"fmt"
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
