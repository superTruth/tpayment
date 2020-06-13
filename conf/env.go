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
}

var configData *ConfigData

func InitConfigData() {
	f, err := os.Open(os.Getenv(EnvFilePath))
	if err != nil {
		panic("can't find config file:" + err.Error())
		return
	}

	fi, err := f.Stat()
	if err != nil {
		panic("config file stat error:" + err.Error())
		return
	}

	dataBytes := make([]byte, fi.Size())
	_, err = f.Read(dataBytes)

	if err != nil {
		panic("read config file error:" + err.Error())
		return
	}

	configData = new(ConfigData)
	err = json.Unmarshal(dataBytes, configData)
	if err != nil {
		panic("config file format error:" + err.Error())
		return
	}
}

func GetConfigData() *ConfigData {
	return configData
}
