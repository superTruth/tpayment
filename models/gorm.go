package models

import (
	"database/sql/driver"
	"tpayment/pkg/gorm_json"
)

type StringArray []string

func (c StringArray) Value() (driver.Value, error) {
	return gorm_json.Value(c)
}

func (c *StringArray) Scan(input interface{}) error {
	return gorm_json.Scan(input, c)
}
