package iso8583

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/iso8583/parse"
	"tpayment/pkg/iso8583/utils"
	"tpayment/pkg/utils/convert_utils"
)

type Message struct {
	fieldConfigs [129]*iso8583Define.FieldConfig // 配置文件
	fieldValue   [129]string                     // 域数据
	headValue    string                          // 头数据(包含TPDU)
}

//
func (message *Message) SetFieldConfig(index int, cfg *iso8583Define.FieldConfig) error {
	if (index < 0) || (index > 128) {
		return errors.New("SetFieldValue out of range->" + strconv.Itoa(index))
	}

	message.fieldConfigs[index] = cfg

	return nil
}

// set field value
func (message *Message) SetFieldValue(index int, value string) error {
	if (index < 0) || (index > 128) {
		return errors.New("SetFieldValue out of range->" + strconv.Itoa(index) + "->" + value)
	}

	if message.fieldConfigs[index] == nil {
		return errors.New("SetFieldValue field->" + strconv.Itoa(index) + "not config")
	}

	message.fieldValue[index] = value
	return nil
}

// get field value
func (message *Message) GetFieldValue(index int) (string, error) {
	if (index < 0) || (index > 128) {
		return "", errors.New("getFieldVaue out of range->" + strconv.Itoa(index))
	}
	return message.fieldValue[index], nil
}

// set header value
func (message *Message) SetHeaderValue(value string) error {
	message.headValue = value
	return nil
}

// get header value
func (message *Message) GetHeaderValue() string {
	return message.headValue
}

// getBitmap value
func (message *Message) GetBitmapValue() (string, error) {
	if message.fieldValue[1] == "" {
		return "", errors.New("didn't generate yet")
	}

	return message.fieldValue[1], nil
}

// set message type
func (message *Message) SetMessageType(value string) error {
	if message.fieldConfigs[0] == nil {
		return errors.New("Message type not config")
	}

	message.fieldValue[0] = value
	return nil
}

// get message type
func (message *Message) GetMessageType() (string, error) {
	if message.fieldValue[0] == "" {
		return "", errors.New("Message type not set")
	}

	return message.fieldValue[0], nil
}

// package message
func (message *Message) FormMessage() ([]byte, error) {

	retSb := bytes.Buffer{}

	// set bitmap
	bitmapStr, err := message.calcBitmap()
	if err != nil {
		return nil, err
	}
	_ = message.SetFieldValue(1, bitmapStr)

	// format message
	for i := 0; i < len(message.fieldConfigs); i++ {
		var tmp []byte
		var err error
		if i == 1 {
			tmp, err = message.formBitmap()
		} else {
			tmp, err = message.formFieldValue(i)
		}

		if err != nil {
			return nil, err
		}
		retSb.Write(tmp)
	}

	return retSb.Bytes(), nil
}

// package message with len and head
func (message *Message) FormMessageX(lenLen int, lenType iso8583Define.FieldValueType) ([]byte, error) {
	retSb := bytes.Buffer{}

	// head
	headBytes := convert_utils.HexString2Bytes(message.headValue)

	// iso8583 message
	iso8583Bytes, err := message.FormMessage()
	if err != nil {
		return nil, err
	}

	// len
	lenBytes := utils.ConvertInt2Bytes(uint64(len(headBytes)+len(iso8583Bytes)), lenType, lenLen)

	// add all
	retSb.Write(lenBytes)
	retSb.Write(headBytes)
	retSb.Write(iso8583Bytes)
	return retSb.Bytes(), nil
}

// parse message
func (message *Message) ParseMessage(buffer []byte, offset int) error {
	offsetCopy := offset
	// message type
	var err error
	message.fieldValue[0], err = message.parseFieldValue(0, buffer, &offsetCopy)
	if err != nil {
		return err
	}
	//fmt.Println("message type->", *message.fieldValue[0]);

	// bitmap
	message.fieldValue[1], err = message.parseBitmapValue(buffer, &offsetCopy)
	if err != nil {
		return err
	}
	//fmt.Println("bitmap->", *message.fieldValue[1]);

	bitmapBytes := convert_utils.HexString2Bytes(message.fieldValue[1])
	for i := 1; i < len(bitmapBytes)*8; i++ {
		enable, err := utils.BitGet(bitmapBytes, i)
		if err != nil {
			return err
		}

		if !enable {
			continue
		}

		message.fieldValue[i+1], err = message.parseFieldValue(i+1, buffer, &offsetCopy)
		if err != nil {
			return err
		}
	}

	return nil
}

//
func (message *Message) formFieldValue(index int) ([]byte, error) {
	if index < 0 || index > 128 {
		return nil, errors.New("formFieldValue index out of range->" + strconv.Itoa(index))
	}

	// didn't config the field
	if message.fieldConfigs[index] == nil {
		return nil, errors.New("formFieldValue didn't config the field->" + strconv.Itoa(index))
	}

	// didn't set value
	if message.fieldValue[index] == "" {
		return make([]byte, 0), nil
	}

	var ret []byte
	switch message.fieldConfigs[index].FieldValueType {
	case iso8583Define.Hex:
		ret = parse.HexFormat(message.fieldConfigs[index], message.fieldValue[index])
	case iso8583Define.Number:
		ret = parse.NumericFormat(message.fieldConfigs[index], message.fieldValue[index])
	case iso8583Define.Alpha:
		ret = parse.AlphaFormat(message.fieldConfigs[index], message.fieldValue[index])
	case iso8583Define.EBCDIC:
		ret = parse.ECDICFormat(message.fieldConfigs[index], message.fieldValue[index])
	}

	return ret, nil
}

//
func (message *Message) formBitmap() ([]byte, error) {
	var ret []byte
	switch message.fieldConfigs[1].FieldValueType {
	case iso8583Define.Hex:
		if len(message.fieldValue[1]) == 16 {
			message.fieldConfigs[1].ValueLen = 8
		} else {
			message.fieldConfigs[1].ValueLen = 16
		}

		ret = parse.HexFormat(message.fieldConfigs[1], message.fieldValue[1])
	case iso8583Define.Number:
		if len(message.fieldValue[1]) == 16 {
			message.fieldConfigs[1].ValueLen = 16
		} else {
			message.fieldConfigs[1].ValueLen = 32
		}
		ret = parse.NumericFormat(message.fieldConfigs[1], message.fieldValue[1])
	case iso8583Define.Alpha:
		if len(message.fieldValue[1]) == 16 {
			message.fieldConfigs[1].ValueLen = 16
		} else {
			message.fieldConfigs[1].ValueLen = 32
		}
		ret = parse.AlphaFormat(message.fieldConfigs[1], message.fieldValue[1])
	case iso8583Define.EBCDIC:
		if len(message.fieldValue[1]) == 16 {
			message.fieldConfigs[1].ValueLen = 16
		} else {
			message.fieldConfigs[1].ValueLen = 32
		}
		ret = parse.ECDICFormat(message.fieldConfigs[1], message.fieldValue[1])
	}

	return ret, nil
}

//
func (message *Message) parseFieldValue(index int, buffer []byte, offset *int) (string, error) {
	if index < 0 || index > 128 {
		return "", errors.New("parseFieldValue index out of range->" + strconv.Itoa(index))
	}

	// didn't config the field
	if message.fieldConfigs[index] == nil {
		return "", errors.New("parseFieldValue didn't config the field->" + strconv.Itoa(index))
	}

	var ret string
	var err error
	switch message.fieldConfigs[index].FieldValueType {
	case iso8583Define.Hex:
		ret, err = parse.HexParse(message.fieldConfigs[index], buffer, offset)
	case iso8583Define.Number:
		ret, err = parse.NumericParse(message.fieldConfigs[index], buffer, offset)
	case iso8583Define.Alpha:
		ret, err = parse.AlphaParse(message.fieldConfigs[index], buffer, offset)
	case iso8583Define.EBCDIC:
		ret, err = parse.ECDICParse(message.fieldConfigs[index], buffer, offset)
	}

	return ret, err
}

//
func (message *Message) parseBitmapValue(buffer []byte, offset *int) (string, error) {
	var ret string
	var err error
	switch message.fieldConfigs[1].FieldValueType {
	case iso8583Define.Hex:
		if buffer[*offset]&0x80 != 0 { // 判断是否是128位数据
			message.fieldConfigs[1].ValueLen = 16
		} else {
			message.fieldConfigs[1].ValueLen = 8
		}
		ret, err = parse.HexParse(message.fieldConfigs[1], buffer, offset)
	case iso8583Define.Number:
		if buffer[*offset] == '1' { // 判断是否是128位数据
			message.fieldConfigs[1].ValueLen = 32
		} else {
			message.fieldConfigs[1].ValueLen = 16
		}
		ret, err = parse.NumericParse(message.fieldConfigs[1], buffer, offset)
	case iso8583Define.Alpha:
		if buffer[*offset] == '1' { // 判断是否是128位数据
			message.fieldConfigs[1].ValueLen = 32
		} else {
			message.fieldConfigs[1].ValueLen = 16
		}
		ret, err = parse.AlphaParse(message.fieldConfigs[1], buffer, offset)
	case iso8583Define.EBCDIC:
		if buffer[*offset] == '1' { // 判断是否是128位数据
			message.fieldConfigs[1].ValueLen = 32
		} else {
			message.fieldConfigs[1].ValueLen = 16
		}
		ret, err = parse.ECDICParse(message.fieldConfigs[1], buffer, offset)
	}

	return ret, err
}

//
func (message *Message) calcBitmap() (string, error) {

	if message.fieldConfigs[1] == nil {
		return "", errors.New("calcBitmap not config")
	}

	// 判断是否有128域
	is128Field := false
	bitmapLen := 8
	for i := 65; i < 129; i++ {
		if message.fieldValue[i] != "" {
			bitmapLen = 16
			is128Field = true
			break
		}
	}

	retBytes := make([]byte, bitmapLen)
	_ = utils.BitSet(retBytes, 0, is128Field)

	for i, fieldValue := range message.fieldValue {
		if (fieldValue == "") || (i < 2) {
			continue
		}
		_ = utils.BitSet(retBytes, i-1, true)
	}

	// 变成string
	ret := ""
	switch message.fieldConfigs[1].FieldValueType {
	case iso8583Define.Alpha:
		ret = convert_utils.Bytes2HexString2Bytes2HexString(retBytes)
	case iso8583Define.Hex:
		ret = convert_utils.Bytes2HexString(retBytes)
	}

	return ret, nil
}

func (message *Message) String() string {
	sb := strings.Builder{}
	sb.WriteString("=============================\n")
	for i, fieldValue := range message.fieldValue {

		if fieldValue == "" {
			continue
		}

		// 不允许打印的就直接跳过
		if message.fieldConfigs[i].Mask {
			sb.WriteString(fmt.Sprintf("field[%d]->********\n", i))
			continue
		}

		if i == 0 {
			sb.WriteString(fmt.Sprintf("message type->%s\n", fieldValue))
			continue
		}

		if i == 1 {
			sb.WriteString(fmt.Sprintf("bitmap->%s\n", fieldValue))
			continue
		}

		sb.WriteString(fmt.Sprintf("field[%d]->%s\n", i, fieldValue))
	}

	sb.WriteString("=============================")

	return sb.String()
}
