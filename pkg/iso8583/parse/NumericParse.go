package parse

import (
	"bytes"
	"errors"
	"fmt"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/iso8583/utils"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/mix_utils"
)

// 打包
func NumericFormat(config *iso8583Define.FieldConfig, value string) []byte {
	if config == nil || value == "" {
		return make([]byte, 0)
	}

	sb := bytes.Buffer{}
	fieldLen := config.ValueLen
	// 如果是不定长数据，需要先把长度写入
	if !config.IsValueLenFix {
		fieldLen = mix_utils.MinInt(config.ValueLen, len(value))
		sb.Write(utils.ConvertInt2Bytes(uint64(fieldLen), config.LenType, config.LenLen))
	}

	// 写入数据
	valueTmp := make([]byte, (fieldLen+1)/2)
	mix_utils.BytesFill(valueTmp, config.PaddingByte) // 数据填充

	valueBytes := []byte(value)
	if len(valueBytes) > fieldLen { // 去掉过长的数据
		valueBytes = mix_utils.BytesArrayCopyArrange(valueBytes, 0, fieldLen)
	}

	// 单字节填充
	if (len(valueBytes) % 2) != 0 { // 单数情况，需要填充
		sbTmp := bytes.Buffer{}

		bytesTmp := make([]byte, 1)
		bytesTmp[0] = config.PaddingByte
		paddingBytes := []byte(convert_utils.Bytes2HexString(bytesTmp))
		switch config.FieldAlignType {
		case iso8583Define.Left:
			sbTmp.Write(valueBytes)
			sbTmp.WriteByte(paddingBytes[1])
		case iso8583Define.Right:
			sbTmp.WriteByte(paddingBytes[0])
			sbTmp.Write(valueBytes)
		}
		valueBytes = sbTmp.Bytes()
	}

	// 处理过后的数据转换成byte数组，再写入
	valueBytes = convert_utils.HexString2Bytes(string(valueBytes))

	switch config.FieldAlignType {
	case iso8583Define.Left:
		mix_utils.BytesArrayCopy(valueBytes, 0, valueTmp, 0, len(valueBytes))
	case iso8583Define.Right:
		mix_utils.BytesArrayCopy(valueBytes, 0, valueTmp, len(valueTmp)-len(valueBytes), len(valueBytes))
	}
	sb.Write(valueTmp)
	return sb.Bytes()
}

// 解包
func NumericParse(config *iso8583Define.FieldConfig, buffer []byte, offset *int) (string, error) {

	bufferLen := len(buffer)
	if *offset < 0 || bufferLen == 0 {
		return "", errors.New("Numeric parse fail")
	}

	valueLen := 0
	if config.IsValueLenFix { // 固定长度解析
		valueLen = config.ValueLen
	} else { // 非定长解析

		llBytes, err := utils.ArrayCopyLenBytes(buffer, offset, config.LenType, config.LenLen) // 内部会有offset指针偏移

		if err != nil {
			return "", err
		}

		valueLen = int(utils.ConvertBytes2Int(llBytes, config.LenType))
	}

	if (*offset + (valueLen+1)/2) > bufferLen {
		return "", fmt.Errorf("Numeric parse fail->value out of range->offset[%d],config.LenLen[%d]", *offset, (valueLen+1)/2)
	}

	ret := convert_utils.Bytes2HexString(mix_utils.BytesArrayCopyArrange(buffer, *offset, (*offset)+(valueLen+1)/2))

	if (valueLen % 2) != 0 {
		switch config.FieldAlignType {
		case iso8583Define.Left:
			ret = string([]byte(ret)[:len(ret)-1])
		case iso8583Define.Right:
			ret = string([]byte(ret)[1:])
		}
	}

	*offset += (valueLen + 1) / 2 // 偏移一下指针

	return ret, nil
}
