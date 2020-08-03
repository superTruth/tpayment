package apkparser

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shogo82148/androidbinary/apk"
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
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return nil, err
	}
	// nolint
	defer f.Close()
	return a.GetApkInfo(filePath)

}

func (a *ApkParser) GetApkInfo(filePath string) (*ApkInfo, error) {
	pkg, err := apk.OpenFile(filePath)
	// nolint
	defer pkg.Close()
	if err != nil {
		fmt.Println("打开APK文件错误:", err)
		return nil, err
	}
	manifest := pkg.Manifest()
	apkInfo := ApkInfo{}
	apkInfo.Package = pkg.PackageName()
	apkInfo.VersionName = manifest.VersionName.MustString()
	apkInfo.VersionCode = int(manifest.VersionCode.MustInt32())
	return &apkInfo, nil
}

//func (a *ApkParser) isFileExist(filePath string) bool {
//	info, err := os.Stat(filePath)
//	if os.IsNotExist(err) {
//		fmt.Println(info)
//		return false
//	}
//	return true
//}

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
