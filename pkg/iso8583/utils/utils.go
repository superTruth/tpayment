package utils

import (
	"errors"
	"strconv"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/mix_utils"
)

func ConvertInt2Bytes(value uint64, bytesType iso8583Define.FieldValueType, len int) []byte {
	switch bytesType {
	case iso8583Define.Alpha:
		return convert_utils.Long2BytesAscii(value, len)
	case iso8583Define.Number:
		return convert_utils.Long2BytesBCD(value, (len+1)/2)
	case iso8583Define.Hex:
		return convert_utils.Long2BytesHex(value, len)
	case iso8583Define.EBCDIC:
		return convert_utils.Long2BytesEBCDIC(value, len)
	}

	return nil
}

func ConvertBytes2Int(value []byte, bytesType iso8583Define.FieldValueType) uint64 {
	switch bytesType {
	case iso8583Define.Alpha:
		return convert_utils.BytesAscii2Long(value, 0, len(value))
	case iso8583Define.Number:
		return convert_utils.BytesBCD2Long(value, 0, len(value))
	case iso8583Define.Hex:
		return convert_utils.BytesHex2Long(value, 0, len(value))
	case iso8583Define.EBCDIC:
		return convert_utils.BytesEBCDIC2Long(value, 0, len(value))
	}

	return 0
}

func ArrayCopyLenBytes(buffer []byte, offset *int, lenType iso8583Define.FieldValueType, length int) ([]byte, error) {
	bufferLen := len(buffer)
	var ret []byte
	switch lenType {
	//case iso8583Define.Hex:
	//	if (*offset + length) > bufferLen { // 提取长度
	//		return nil, errors.New("ArrayCopyLenBytes Hex out of range")
	//	}
	//	ret = mix_utils.BytesArrayCopyArrange(buffer, *offset, (*offset)+length)
	//	*offset += length
	//	break
	case iso8583Define.Number:
		if (*offset + (length+1)/2) > bufferLen { // 提取长度
			return nil, errors.New("ArrayCopyLenBytes Number out of range")
		}
		ret = mix_utils.BytesArrayCopyArrange(buffer, *offset, (*offset)+(length+1)/2)
		*offset += (length + 1) / 2
	case iso8583Define.Alpha, iso8583Define.EBCDIC, iso8583Define.Hex:
		if (*offset + length) > bufferLen { // 提取长度
			return nil, errors.New("ArrayCopyLenBytes Alpha out of range")
		}
		ret = mix_utils.BytesArrayCopyArrange(buffer, *offset, (*offset)+length)
		*offset += length
	}

	return ret, nil
}

func BitSet(value []byte, index int, enable bool) error {
	if index < 0 || index > len(value)*8 {
		return errors.New("BitSet out of range index->" + strconv.Itoa(index))
	}

	destByteIndex := index / 8
	destBitIndex := index % 8

	if enable {
		value[destByteIndex] |= byte(1 << uint(7-destBitIndex))
	} else {
		value[destByteIndex] &= ^byte(1 << uint(7-destBitIndex))
	}

	return nil
}

func BitGet(value []byte, index int) (bool, error) {
	if index < 0 || index > len(value)*8 {
		return false, errors.New("BitGet out of range index->" + strconv.Itoa(index))
	}

	destByteIndex := index / 8
	destBitIndex := index % 8

	ret := value[destByteIndex]&byte(1<<uint(7-destBitIndex)) != 0

	return ret, nil
}
