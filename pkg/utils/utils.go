package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func Body2Json(body io.Reader, destBean interface{}) error {

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &destBean)
}

func Struct2Map(src interface{}) []map[string]interface{} {
	srcByte,_ := json.Marshal(src)

	var ret []map[string]interface{}

	_ = json.Unmarshal(srcByte, &ret)

	return ret
}

// 只保留剩余的key数据，删除多余的
func FilterMapData(src map[string]interface{}, keys map[string]bool) {
	for k,_ := range src {
		if !keys[k] {
			delete(src, k)
		}
	}
}
