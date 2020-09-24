package apkparser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestApkParser(t *testing.T) {
	//{"version_name":"1.0.18_0-1-6","version_code":5,"package":"com.octopuscards.octopusframework"}
	//{"version_name":"1.0.18_0-1-5","version_code":5,"package":"com.octopuscards.octopusframework"}
	filePath := "/Users/truth/Downloads/horizon_mdm.apk"

	_, err := os.OpenFile(filePath, os.O_RDONLY, 0)

	if err != nil {
		fmt.Println("err->", err.Error())
		return
	}

	apkParser := ApkParser{}
	apkInfo, err := apkParser.GetApkInfo(filePath)
	if err != nil {
		fmt.Println("err->", err.Error())
		return
	}

	apkByte, _ := json.Marshal(*apkInfo)
	fmt.Println("apk Info", string(apkByte))
}
