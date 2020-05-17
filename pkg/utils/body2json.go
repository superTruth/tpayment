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

