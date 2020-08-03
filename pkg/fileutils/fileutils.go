package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//
func CreateLocalPath(destPath string) error {
	dirPath, _, _ := SeparateFilePath(destPath)

	if b, _ := IsPathExists(dirPath); !b {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// 判断本地文件是否存在
func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileList(path string) []string {
	var ret []string
	fs, _ := ioutil.ReadDir(path)

	for _, file := range fs {
		ret = append(ret, file.Name())
	}

	return ret
}

func SeparateFilePath(destPath string) (string, string, error) {

	fileName := path.Base(destPath)
	fmt.Println("SeparateFilePath->", destPath, ",", fileName)

	dirPath := destPath[0 : len(destPath)-len(fileName)]

	return dirPath, fileName, nil
}

func DeleteFile(localPath string) error {
	return os.Remove(localPath)
}
