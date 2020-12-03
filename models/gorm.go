package models

import (
	"database/sql/driver"
	"strconv"
	"strings"
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

func (c *StringArray) String() string {
	sb := strings.Builder{}

	for i := 0; i < len(*c); i++ {
		if sb.Len() != 0 {
			sb.WriteString(",")
		}
		sb.WriteString((*c)[i])
	}

	return sb.String()
}

// int类型
type IntArray []uint64

func (c IntArray) Value() (driver.Value, error) {
	return gorm_json.Value(c)
}
func (c *IntArray) Scan(input interface{}) error {
	return gorm_json.Scan(input, c)
}

func (c *IntArray) Change2UintArray() []uint64 {
	tmp := make([]uint64, len(*c))

	for i := 0; i < len(tmp); i++ {
		tmp[i] = (*c)[i]
	}

	return tmp
}

func (c *IntArray) String() string {
	sb := strings.Builder{}

	for i := 0; i < len(*c); i++ {
		if sb.Len() != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatUint(uint64((*c)[i]), 10))
	}

	return sb.String()
}
