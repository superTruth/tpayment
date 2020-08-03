package conf

import (
	"encoding/json"
	"os"
)

const (
	EnvFilePath = "ENV_FILE_PATH"
)

type ConfigData struct {
	Domain string `json:"domain"`

	EmailHost        string `json:"email_host"`
	EmailHostPort    int    `json:"email_host_port"`
	EmailUserAccount string `json:"email_user_account"`
	EmailUserPwd     string `json:"email_user_pwd"`
	EmailUserName    string `json:"email_user_name"`

	DBAccount string `json:"db_account"`

	S3Region string `json:"s3_region"`
	S3Key    string `json:"s3_key"`
	S3Secret string `json:"s3_secret"`
	S3Bucket string `json:"s3_bucket"`
}

var configData *ConfigData

func InitConfigData() {
	f, err := os.Open(os.Getenv(EnvFilePath))
	if err != nil {
		panic("can't find config fileutils:" + err.Error())
	}

	fi, err := f.Stat()
	if err != nil {
		panic("config fileutils stat error:" + err.Error())
	}

	dataBytes := make([]byte, fi.Size())
	_, err = f.Read(dataBytes)

	if err != nil {
		panic("read config fileutils error:" + err.Error())
	}

	configData = new(ConfigData)
	err = json.Unmarshal(dataBytes, configData)
	if err != nil {
		panic("config fileutils format error:" + err.Error())
	}
}

func GetConfigData() *ConfigData {
	return configData
}
