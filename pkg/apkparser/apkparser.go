package apkparser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/avast/apkparser"

	"github.com/google/uuid"
)

type ApkParser struct {
	Url string // APK url
}

var localPath = "./apkdecodedir/"

type ApkInfo struct {
	VersionName string `json:"version_name"`
	VersionCode int    `json:"version_code"`
	Package     string `json:"package"`
}

// Download APK and
func (a *ApkParser) DownloadApkInfo() (*ApkInfo, error) {

	// 检查LocalPath是否创建
	basePath, err := checkAndCreateLocalPath()
	if err != nil {
		return nil, err
	}

	filename := uuid.New().String() //path.Base(uri.Path)		// 获取文件名字
	fmt.Println("[*] Filename " + filename)
	filePath := basePath + "/" + filename
	defer os.Remove(filePath) // 程序执行完毕删除apk文件

	// 文件不存在，进行下载
	res, err := http.Get(a.Url)
	if err != nil {
		fmt.Println("[*] Get url failed:" + a.Url)
		return nil, err
	}

	//创建下载存放apk
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("[*] Create temp fileutils failed:", err)
		return nil, err
	}

	// nolint
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return nil, err
	}

	return a.GetApkInfo(filePath)
}

type ManifestBean struct {
	XMLName     xml.Name `xml:"manifest"`
	Package     string   `xml:"package,attr"`
	VersionCode int      `xml:"versionCode,attr"`
	VersionName string   `xml:"versionName,attr"`
}

func (a *ApkParser) GetApkInfo(filePath string) (*ApkInfo, error) {

	buf := new(bytes.Buffer)
	enc := xml.NewEncoder(buf)
	enc.Indent("", "\t")
	zipErr, resErr, manErr := apkparser.ParseApk(filePath, enc)
	if zipErr != nil {
		fmt.Println("Failed to open the APK: ", zipErr.Error())
		return nil, errors.New("zip err->" + zipErr.Error())
	}

	if resErr != nil {
		fmt.Println("Failed to parse resources:", resErr.Error())
		return nil, errors.New("resErr err->" + resErr.Error())
	}
	if manErr != nil {
		fmt.Println("Failed to parse AndroidManifest.xml:", manErr.Error())
		return nil, errors.New("manErr err->" + resErr.Error())
	}

	manifestBean := new(ManifestBean)
	err := xml.Unmarshal(buf.Bytes(), manifestBean)
	if err != nil {
		return nil, err
	}

	apkInfo := ApkInfo{}
	apkInfo.Package = manifestBean.Package
	apkInfo.VersionName = manifestBean.VersionName
	apkInfo.VersionCode = manifestBean.VersionCode
	return &apkInfo, nil
}

// 创建本地缓存文件夹
func checkAndCreateLocalPath() (string, error) {
	if b, _ := isPathExists(localPath); !b {
		if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	//
	currentDate := time.Now().Format("20060102")
	currentDateInt, _ := strconv.Atoi(currentDate)

	fileList := getFileList(localPath)
	for _, fileName := range fileList { // 遍历文件夹，把超过2天的文件夹删除掉，以免磁盘塞满
		date, _ := strconv.Atoi(fileName)
		if (currentDateInt - date) > 1 {
			os.Remove(localPath + fileName)
		}
	}

	// 如果文件夹已经存在，则不需要再次创建
	if b, _ := isPathExists(localPath + currentDate); !b {
		return localPath + currentDate, os.MkdirAll(localPath+currentDate, os.ModePerm)
	}

	return localPath + currentDate, nil
}

// 判断本地文件夹是否存在
func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getFileList(path string) []string {
	var ret []string
	fs, _ := ioutil.ReadDir(path)

	for _, file := range fs {
		ret = append(ret, file.Name())
	}

	return ret
}
