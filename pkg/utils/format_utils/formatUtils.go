package format_utils

import "tpayment/pkg/utils/mix_utils"

/**
 * append arrays
 *
 * @param data
 * @param destLen
 * @param appendToFront
 * @param append
 * @return
 */
func AppendBytes(value []byte, destLen int, appendToFront bool, appendByte byte) []byte {
	ret := make([]byte, destLen)
	mix_utils.BytesFill(ret, appendByte)

	if value == nil {
		return ret
	}

	if len(value) == destLen {
		return value
	}

	if len(value) > destLen { //  need cut string instead of append
		if appendToFront {
			mix_utils.BytesArrayCopy(value, len(value)-destLen, ret, 0, destLen)
			return ret
		}
		mix_utils.BytesArrayCopy(value, 0, ret, 0, destLen)
		return ret
	}

	if appendToFront {
		mix_utils.BytesArrayCopy(value, 0, ret, destLen-len(value), len(value))
		return ret
	}
	mix_utils.BytesArrayCopy(value, 0, ret, 0, len(value))

	return ret
}

/**
* append arrays
*
* @param data
* @param destLen
* @param appendToFront
* @param append
* @return
 */
func AppendString(value string, destLen int, appendToFront bool, appendByte byte) string {
	valueByte := []byte(value)
	valueByte = AppendBytes(valueByte, destLen, appendToFront, appendByte)
	return string(valueByte)
}

func DeleteAppendBytes(value []byte, deleteToFront bool, deleteByte byte) []byte {
	if value == nil {
		return nil
	}

	index := 0
	if deleteToFront {
		for ; index < len(value); index++ {
			if value[index] != deleteByte {
				break
			}
		}
		if index == len(value) {
			return make([]byte, 0)
		}
		return mix_utils.BytesArrayCopyArrange(value, index, len(value))
	} else {
		for index = (len(value) - 1); index >= 0; index-- {
			if value[index] != deleteByte {
				break
			}
		}
		if index == 0 {
			return make([]byte, 0)
		}
		return mix_utils.BytesArrayCopyArrange(value, 0, index+1)
	}

}

func DeleteAppendString(value string, deleteToFront bool, deleteByte byte) string {
	valueBytes := []byte(value)

	valueBytes = DeleteAppendBytes(valueBytes, deleteToFront, deleteByte)

	return string(valueBytes)
}
