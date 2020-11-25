package convert_utils

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"tpayment/pkg/utils/format_utils"
	"tpayment/pkg/utils/mix_utils"
)

var hexDigits = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}

/**
* byteArr -> hexString
* <p>exp：</p>
* Bytes2HexString(new byte[] { 0, (byte) 0xa8 }, ", ") returns 00, A8
*
* @param bytes source data
* @return upcase result string
 */
func Bytes2HexStringX(value []byte, split string) string {
	if (value == nil) || len(value) == 0 {
		return ""
	}
	sb := bytes.Buffer{}
	for i, c := range value {
		sb.WriteByte(hexDigits[(c>>4)&0x0F])
		sb.WriteByte(hexDigits[c&0x0F])
		if i != (len(value) - 1) {
			sb.WriteString(split)
		}
	}

	return sb.String()
}

/**
 * byteArr -> hexString
 * <p>exp：</p>
 * bytes2HexString(new byte[] { 0, (byte) 0xa8 }) returns 00A8
 * @param bytes
 * @return
 */
func Bytes2HexString(value []byte) string {

	return Bytes2HexStringX(value, "")
}

/**
* hexString -> byteArr
* <p>exp：</p>
* hexString2Bytes("00A8") returns { 0, (byte) 0xA8 }
*
* @param hexString source data
* @return result byte array
 */
func HexString2Bytes(hexString string) []byte {
	tmp := FormatHexString(hexString)
	if len(tmp) == 0 {
		return nil
	}

	tmpLen := len(tmp)

	hexBytes := []byte(tmp)

	ret := make([]byte, tmpLen>>1)

	for i := 0; i < tmpLen; i += 2 {
		ret[i>>1] = (byte)(hex2Dec(hexBytes[i])<<4 | hex2Dec(hexBytes[i+1]))
	}

	return ret
}

/**
* long -> hex bytes
* <p>exp：</p>
* long2BytesHex(1025, 3) returns {00, 04, 01}
*
* @param value
* @param len return buffer len
* @return
 */
func Long2BytesHex(value uint64, length int) []byte {
	ret := make([]byte, length)
	mix_utils.BytesFill(ret, 0)

	tmp := value
	for i := 0; i < len(ret); i++ {
		ret[len(ret)-1-i] = (byte)(tmp & 0xFF)
		tmp >>= 8
	}
	return ret
}

func Long2BytesEBCDIC(value uint64, length int) []byte {
	valueStr := strconv.FormatUint(value, 10)
	return ASCII2ECDIC(String2PString(format_utils.AppendString(valueStr, length, true, '0')))
}

func BytesEBCDIC2Long(value []byte, offset int, length int) uint64 {
	var ret uint64 = 0

	if (value == nil) || (len(value) <= 0) {
		return ret
	}

	if (offset + length) > len(value) {
		return ret
	}

	for i := 0; i < length; i++ {
		ret *= 10
		ret += uint64(value[i] - 0xF0)
	}

	return ret
}

/**
* hex bytes -> long
* <p>exp：</p>
* bytesHex2Long({ 0x11, 0xA8, 0x23 }, 1, 2) returns 43043
*
* @param bytes
* @param offset
* @param len
* @return
 */
func BytesHex2Long(value []byte, offset int, length int) uint64 {
	var ret uint64 = 0

	if (value == nil) || (len(value) == 0) {
		return ret
	}

	if (offset + length) > len(value) {
		return ret
	}

	for i := 0; i < length; i++ {
		ret <<= 8
		ret |= uint64(value[offset+i] & 0xFF)
	}

	return ret
}

/**
* long -> BCD bytes
* <p>exp：</p>
* long2BytesBCD(1223, 3) returns { 0x00, 0x12, 0x23 }
*
* @param value
* @param len return buffer len
* @return
 */
func Long2BytesBCD(value uint64, length int) []byte {
	ret := make([]byte, length)
	mix_utils.BytesFill(ret, 0)

	tmp := value

	for i := 0; i < len(ret); i++ {
		ret[len(ret)-1-i] = int2BCD((int)(tmp % 100))
		tmp /= 100
	}

	return ret
}

/**
* BCD bytes -> long
* <p>exp：</p>
* bytesBCD2Long({ 0x11, 0x12, 0x23 }, 1, 2) returns 1223
*
* @param bytes
* @param offset
* @param len
* @return
 */
func BytesBCD2Long(value []byte, offset int, length int) uint64 {
	var ret uint64 = 0

	if (value == nil) || (len(value) <= 0) {
		return ret
	}

	if (offset + length) > len(value) {
		return ret
	}

	for i := 0; i < length; i++ {
		ret *= 100
		ret += uint64(bcd2Int(value[offset+i]))
	}

	return ret
}

//
func String2PString(value string) *string {
	return &value
}

func Bool2PBool(value bool) *bool {
	return &value
}

func Int2PInt(value int) *int {
	return &value
}

func Uint642PUint64(value uint64) *uint64 {
	return &value
}

func Long2BytesAscii(value uint64, length int) []byte {
	return format_utils.AppendBytes([]byte(strconv.FormatUint(value, 10)), length, true, '0')
}

func BytesAscii2Long(value []byte, offset int, length int) uint64 {
	tmp := string(mix_utils.BytesArrayCopyArrange(value, offset, offset+length))
	ret, err := strconv.ParseUint(tmp, 10, 64)

	if err != nil {
		return 0
	}
	return ret
}

func Bytes2HexString2Bytes2HexString(value []byte) string {
	return Bytes2HexString([]byte(Bytes2HexString(value)))
}

func hex2Dec(hexChar byte) int {
	if (hexChar >= '0') && (hexChar <= '9') {
		return int(hexChar) - int('0')
	} else {
		return int(hexChar) - int('A') + 10
	}
}

func int2BCD(value int) byte {
	if value > 99 {
		return 0
	}
	return (byte)(value/10*16 + value%10)
}

func bcd2Int(b byte) int {
	tmpb := int(b & 0xFF)
	return tmpb/16*10 + tmpb%16
}

/**
* delete character in hex string what's not hex format
* <p>exp：</p>
* formatHexString("0x00, 0xA8") returns "00A8"
*
* @param hexString
* @return hex string
 */
func FormatHexString(hexString string) string {
	if len(hexString) == 0 {
		return ""
	}

	ret := strings.ToUpper(hexString)

	formatCorrect := true
	for c := range ret {
		if ((c >= '0') && (c <= '9')) || ((c >= 'A') && (c <= 'F')) {
			continue
		}
		formatCorrect = false
		break
	}

	if formatCorrect {
		if len(ret)%2 != 0 {
			return "0" + ret
		}
		return ret
	}

	ret = strings.Replace(ret, "0X", "", -1)

	retChars := []byte(ret)

	sb := bytes.Buffer{}
	for i := 0; i < len(retChars); i++ {
		c := retChars[i]
		if ((c >= '0') && (c <= '9')) || ((c >= 'A') && (c <= 'F')) {
			sb.WriteByte(byte(c))
			continue
		}
	}

	if sb.Len()%2 != 0 {
		sbTmp := bytes.Buffer{}
		sbTmp.WriteByte('0')
		sbTmp.Write(sb.Bytes())
		sb = sbTmp
	}

	return sb.String()
}

var ASCII2EBCDIC_DATA = []byte{'\u0000', '\u0001', '\u0002', '\u0003', '\u0004', '\u0005', '\u0006', '\u0007', '\b', '\t', '\n', '\u000b', '\f', '\r', '\u000e', '\u000f', '\u0010', '\u0011', '\u0012', '\u0013', '\u0014', '\u0015', '\u0016', '\u0017', '\u0018', '\u0019', '\u001a', '\u001b', '\u001c', '\u001d', '\u001e', '\u001f', '@', 'Z', '\u007f', '{', '[', 'l', 'P', '}', 'M', ']', '\\', 'N', 'k', '`', 'K', 'a', 'ð', 'ñ', 'ò', 'ó', 'ô', 'õ', 'ö', '÷', 'ø', 'ù', 'z', '^', 'L', '~', 'n', 'o', '|', 'Á', 'Â', 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'â', 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', '\u00ad', 'à', '½', '_', 'm', '}', '\u0081', '\u0082', '\u0083', '\u0084', '\u0085', '\u0086', '\u0087', '\u0088', '\u0089', '\u0091', '\u0092', '\u0093', '\u0094', '\u0095', '\u0096', '\u0097', '\u0098', '\u0099', '¢', '£', '¤', '¥', '¦', '§', '¨', '©', 'À', 'j', 'Ð', '¡', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K', 'K'}
var EBCDIC2ASCII_DATA = []byte{'\u0000', '\u0001', '\u0002', '\u0003', '\u0004', '\u0005', '\u0006', '\u0007', '\b', '\t', '\n', '\u000b', '\f', '\r', '\u000e', '\u000f', '\u0010', '\u0011', '\u0012', '\u0013', '\u0014', '\u0015', '\u0016', '\u0017', '\u0018', '\u0019', '\u001a', '\u001b', '\u001c', '\u001d', '\u001e', '\u001f', ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/', '.', '.', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '.', '?', ' ', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '<', '(', '+', '|', '&', '.', '.', '.', '.', '.', '.', '.', '.', '.', '!', '$', '*', ')', ';', '^', '-', '/', '.', '.', '.', '.', '.', '.', '.', '.', '|', ',', '%', '_', '>', '?', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', ':', '#', '@', '\'', '=', '"', '.', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', '.', '.', '.', '.', '.', '.', '.', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', '.', '.', '.', '.', '.', '.', '.', '~', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '.', '.', '.', '[', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', ']', '.', '.', '{', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', '.', '.', '.', '.', '.', '.', '}', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', '.', '.', '.', '.', '.', '.', '\\', '.', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '.', '.', '.', '.', '.', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '.', '.', '.', '.', '.'}

func ASCII2ECDIC(value *string) []byte {

	if (value == nil) || (len(*value) == 0) {
		return make([]byte, 0)
	}

	valueByte := []byte(*value)

	for i, v := range valueByte {
		valueByte[i] = ASCII2EBCDIC_DATA[v]
	}

	return valueByte
}

func ECDIC2ASCII(value []byte) *string {
	if (value == nil) || (len(value) == 0) {
		return String2PString("")
	}

	for i, v := range value {
		value[i] = EBCDIC2ASCII_DATA[v]
	}

	return String2PString(string(value))
}

func JsonFormat(v interface{}) (*string, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return String2PString(string(jsonBytes)), nil
}
