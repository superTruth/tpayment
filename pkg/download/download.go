package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"tpayment/pkg/fileutils"
)

func Download(sourceUrl, destUrl string) error {
	err := fileutils.CreateLocalPath(destUrl)

	if err != nil {
		return err
	}

	res, err := http.Get(sourceUrl)
	if err != nil {
		fmt.Println("[*] Get url failed:" + sourceUrl + "," + err.Error())
		return err
	}

	//创建下载存放apk
	f, err := os.Create(destUrl)
	if err != nil {
		fmt.Println("[*] Create temp fileutils failed:", err)
		return err
	}
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	// nolint
	defer f.Close()

	return nil
}
