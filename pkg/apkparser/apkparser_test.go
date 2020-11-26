package apkparser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/avast/apkparser"
)

func TestApkParser(t *testing.T) {
	//{"version_name":"1.0.18_0-1-6","version_code":5,"package":"com.octopuscards.octopusframework"}
	//{"version_name":"1.0.18_0-1-5","version_code":5,"package":"com.octopuscards.octopusframework"}
	filePath := "/Users/truth/Downloads/SeitoOnlineOrderScan_V1.0.3_release.apk" //"/Users/truth/Downloads/horizon_mdm.apk"

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

func TestApkParser2(t *testing.T) {
	filePath := "/Users/truth/Downloads/horizon_demo.apk"

	buf := new(bytes.Buffer)
	enc := xml.NewEncoder(buf)
	enc.Indent("", "\t")
	zipErr, resErr, manErr := apkparser.ParseApk(filePath, enc)
	if zipErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to open the APK: %s", zipErr.Error())
		os.Exit(1)
		return
	}

	if resErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse resources: %s", resErr.Error())
	}
	if manErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse AndroidManifest.xml: %s", manErr.Error())
		os.Exit(1)
		return
	}

	fmt.Println("ret->", buf.String())
	manifestBean := new(ManifestBean)
	err := xml.Unmarshal(buf.Bytes(), manifestBean)
	if err != nil {
		return
	}

	fmt.Println("buf->", manifestBean)
	//fmt.Println()
}
