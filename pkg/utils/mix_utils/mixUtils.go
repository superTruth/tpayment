package mix_utils

import (
	"bytes"
	"errors"
	"math/rand"
	"strings"
)

func BytesArrayCopy(src []byte, srcOffset int, dest []byte, destOffset int, length int) {
	for i := 0; i < length; i++ {
		dest[destOffset+i] = src[srcOffset+i]
	}
}

func BytesArrayCopyArrange(src []byte, fromIndex int, endIndex int) []byte {
	ret := make([]byte, endIndex-fromIndex)
	for i := 0; i < len(ret); i++ {
		ret[i] = src[i+fromIndex]
	}

	return ret
}

func BytesFill(src []byte, fillValue byte) {
	for i := 0; i < len(src); i++ {
		src[i] = fillValue
	}
}

func MinInt(value1 int, value2 int) int {
	if value1 > value2 {
		return value2
	}
	return value1
}

func MaxInt(value1 int, value2 int) int {
	if value1 > value2 {
		return value1
	}
	return value2
}

func Xor(value1 []byte, value2 []byte) ([]byte, error) {
	if value1 == nil || value2 == nil {
		return nil, errors.New("data can't nil")
	}

	if len(value1) != len(value2) {
		return nil, errors.New("two data len not equal")
	}

	ret := make([]byte, len(value1))
	for i := 0; i < len(value1); i++ {
		ret[i] = value1[i] ^ value2[i]
	}

	return ret, nil
}

func SubString(src string, start int, end int) string {
	ret := string([]byte(src)[start:end])
	return ret
}

func SubStringStart(src *string, start int) *string {
	ret := string([]byte(*src)[start:])
	return &ret
}

func Compare(data1 []byte, date2 []byte) bool {
	if data1 == nil || date2 == nil {
		return false
	}

	if len(data1) != len(date2) {
		return false
	}

	for i := 0; i < len(data1); i++ {
		if data1[i] != date2[i] {
			return false
		}
	}
	return true
}

func MergeBytesArray(data1 []byte, data2 []byte) []byte {
	retSb := bytes.Buffer{}

	retSb.Write(data1)
	retSb.Write(data2)

	return retSb.Bytes()
}

// 检查变量是否为整数
func CheckNum(value *string) bool {
	if value == nil {
		return false
	}

	valueBytes := []byte(*value)
	for i := 0; i < len(valueBytes); i++ {
		if valueBytes[i] < '0' || valueBytes[i] > '9' {
			return false
		}
	}

	return true
}

// 检查变量是否为16进制
func CheckHex(value *string) bool {
	if value == nil {
		return false
	}

	valueBytes := []byte(strings.ToUpper(*value))
	for i := 0; i < len(valueBytes); i++ {
		if (valueBytes[i] >= '0' && valueBytes[i] <= '9') || (valueBytes[i] >= 'A' && valueBytes[i] <= 'F') {
			continue
		}
		return false
	}

	return true
}

// 检查是否为空
func CheckEmptyString(value *string) bool {
	if (value == nil) || len(*value) == 0 {
		return true
	}
	return false
}
func CheckEmptyByte(value []byte) bool {
	if (value == nil) || len(value) == 0 {
		return true
	}
	return false
}

// 产生随机数组
func RadomByteArray(value []byte) {
	for i := 0; i < len(value); i++ {
		value[i] = byte(rand.Intn(255))
	}
}
