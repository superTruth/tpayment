package gorm_json

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

func Value(l interface{}) (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func Scan(input interface{}, l interface{}) (err error) {
	switch value := input.(type) {
	case string:
		err = json.Unmarshal([]byte(value), l)
	case []byte:
		err = json.Unmarshal(value, l)
	default:
		err = errors.New("not supported type")
	}
	return
}
