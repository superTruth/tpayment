package parse

import (
	"bytes"
	"errors"
	"fmt"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/iso8583/utils"
	"tpayment/pkg/utils/mix_utils"
)

// 打包
func AlphaFormat(config *iso8583Define.FieldConfig, value string) []byte {
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
	valueTmp := make([]byte, fieldLen)
	mix_utils.BytesFill(valueTmp, config.PaddingByte) // 数据填充

	valueBytes := []byte(value)
	if len(valueBytes) > fieldLen { // 去掉过长的数据
		valueBytes = mix_utils.BytesArrayCopyArrange(valueBytes, 0, fieldLen)
	}

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
func AlphaParse(config *iso8583Define.FieldConfig, buffer []byte, offset *int) (string, error) {
	bufferLen := len(buffer)

	if *offset < 0 || bufferLen == 0 {
		return "", errors.New("alpha parse fail")
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

	if (*offset + valueLen) > bufferLen {
		return "", fmt.Errorf("alpha parse fail->value out of range->offset[%d],len[%d]", *offset, valueLen)
	}

	ret := string(mix_utils.BytesArrayCopyArrange(buffer, *offset, (*offset)+valueLen))
	*offset += valueLen // 偏移一下指针

	return ret, nil
}
