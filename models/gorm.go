package models

import (
	"database/sql/driver"
	"tpayment/pkg/gorm_json"
)

// string类型
type StringArray []string

func (c StringArray) Value() (driver.Value, error) {
	return gorm_json.Value(c)
}

func (c *StringArray) Scan(input interface{}) error {
	return gorm_json.Scan(input, c)
}

// int类型
type IntArray []uint

func (c IntArray) Value() (driver.Value, error) {
	return gorm_json.Value(c)
}
func (c *IntArray) Scan(input interface{}) error {
	return gorm_json.Scan(input, c)
}

func (c *IntArray) Change2UintArray() []uint {
	tmp := make([]uint, len(*c))

	for i := 0; i < len(tmp); i++ {
		tmp[i] = (*c)[i]
	}

	return tmp
}
