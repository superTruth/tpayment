package iso8583

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/utils/convert_utils"
)

// 解析后配置
type Factory struct {
	Configs [129]*iso8583Define.FieldConfig
}

// 创建一个数据对象
func (f *Factory) GenerateNewMessage() *Message {
	msg := Message{}

	for i := 0; i < len(f.Configs); i++ {
		msg.fieldConfigs[i] = f.Configs[i]
	}

	return &msg
}

// 解析配置文件
func CreateConfigFactory(configFilePath string) (*Factory, error) {
	xmlConfig, err := parseXMLFile(configFilePath)
	if err != nil {
		return nil, err
	}

	fieldConfigs, err := parseXMLConfig(xmlConfig)
	if err != nil {
		return nil, err
	}

	factory := Factory{fieldConfigs}

	return &factory, nil
}

// *********************************配置文件反序列化*********************************
type iso8583XMLConfig struct {
	XMLName     xml.Name                `xml:"cn8583-config"`
	FieldConfig []iso8583XMLConfigField `xml:"field"`
}

type iso8583XMLConfigField struct {
	Index   string `xml:"index,attr"`
	Type    string `xml:"type,attr"`
	Align   string `xml:"align,attr"`
	LenType string `xml:"lenType,attr"`
	Padding string `xml:"padding,attr"`
	Mask    string `xml:"mask,attr"`
}

func parseXMLFile(configFilePath string) (*iso8583XMLConfig, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("open file err->", err.Error())
		return nil, err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read file err->", err.Error())
		return nil, err
	}

	config := iso8583XMLConfig{}
	err = xml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("read file err->", err.Error())
		return nil, err
	}
	return &config, nil
}

// *********************************解析配置文件*********************************
func parseXMLConfig(xmlConfig *iso8583XMLConfig) ([129]*iso8583Define.FieldConfig, error) {
	ret := [129]*iso8583Define.FieldConfig{}
	for _, fieldXMLConfig := range xmlConfig.FieldConfig {
		// index parse
		if len(fieldXMLConfig.Index) == 0 {
			return ret, errors.New("Index can't empty")
		}
		index, err := strconv.Atoi(fieldXMLConfig.Index)
		if (err != nil) || (index < 0 || index > 128) {
			return ret, errors.New("Index parse fail->" + fieldXMLConfig.Index)
		}

		ret[index] = &iso8583Define.FieldConfig{}

		// type parse
		ret[index].FieldValueType, ret[index].IsValueLenFix, ret[index].ValueLen, ret[index].LenLen, err = parseTypeAndLen(fieldXMLConfig.Type)
		if err != nil {
			return ret, errors.New(err.Error() + "->" + strconv.Itoa(index))
		}

		// align type parse
		ret[index].FieldAlignType, err = parseAlign(fieldXMLConfig.Align)
		if err != nil {
			return ret, errors.New(err.Error() + "->" + strconv.Itoa(index))
		}

		// len type parse
		ret[index].LenType, err = parseLenType(fieldXMLConfig.LenType)
		if err != nil {
			return ret, errors.New(err.Error() + "->" + strconv.Itoa(index))
		}

		// padding parse
		ret[index].PaddingByte, err = parsePadding(fieldXMLConfig.Padding)
		if err != nil {
			return ret, errors.New(err.Error() + "->" + strconv.Itoa(index))
		}

		// mask parse
		ret[index].Mask = parseMask(fieldXMLConfig.Mask)
	}

	return ret, nil
}

func parseTypeAndLen(value string) (valueType iso8583Define.FieldValueType, isFixLen bool, length int, lenLength int, err error) {
	valueBytes := []byte(strings.ToUpper(value))
	if len(valueBytes) < 2 {
		err = errors.New("field type err")
		return
	}

	// 匹配最后数字
	numStartFlag := false
	sb := bytes.Buffer{}
	for i := 0; i < len(valueBytes); i++ {
		// 查看第一位是否正确
		if i == 0 {
			if valueBytes[0] != 'A' && valueBytes[0] != 'N' && valueBytes[0] != 'H' && valueBytes[0] != 'E' {
				err = errors.New("field type must start with A/N/H/E")
				return
			}
			continue
		}

		// 拷贝最后的数字
		if valueBytes[i] >= '0' && valueBytes[i] <= '9' {
			numStartFlag = true
			sb.WriteByte(valueBytes[i])
			continue
		}

		// 如果已经遇到了数字不是连续的，则报错
		if numStartFlag {
			err = errors.New("field type number is not correct")
			return
		}

		// 包含有非法字符
		if valueBytes[i] != '.' {
			err = errors.New("field type contain unknow char->" + string(valueBytes[i]))
			return
		}
	}
	// field type
	switch valueBytes[0] {
	case 'A':
		valueType = iso8583Define.Alpha
	case 'N':
		valueType = iso8583Define.Number
	case 'H':
		valueType = iso8583Define.Hex
	case 'E':
		valueType = iso8583Define.EBCDIC
	}
	// length
	length, _ = strconv.Atoi(sb.String())

	// 计算小数点的位数
	lenLength = strings.Count(value, ".") // 所有长度 - 数字长度 - 开始位 - 一个点
	if lenLength <= 0 {                   // fix length field
		isFixLen = true
		lenLength = 0
	} else {
		isFixLen = false
	}

	return
}

func parseAlign(value string) (align iso8583Define.FieldAlignType, err error) {
	switch value {
	case "LEFT":
		align = iso8583Define.Left
		return
	case "RIGHT":
		align = iso8583Define.Right
		return
	default:
		err = errors.New("Unknow Align type->" + value)
		return
	}
}

func parseLenType(value string) (lenType iso8583Define.FieldValueType, err error) {
	switch value {
	case "A":
		lenType = iso8583Define.Alpha
	case "N":
		lenType = iso8583Define.Number
	case "H":
		lenType = iso8583Define.Hex
	case "E":
		lenType = iso8583Define.EBCDIC
	default:
		err = errors.New("Unknow Len type->" + value)
	}
	return
}

func parsePadding(value string) (paddingByte byte, err error) {
	valueBytes := []byte(strings.ToUpper(value))
	if len(valueBytes) != 2 {
		err = errors.New("padding byte error->" + value)
		return
	}

	for i := 0; i < len(valueBytes); i++ {
		if ((valueBytes[i] >= '0') && (valueBytes[i] <= '9')) || ((valueBytes[i] >= 'A') && (valueBytes[i] <= 'F')) {

		} else {
			err = errors.New("padding byte not hex->" + value)
			break
		}
	}

	paddingByte = convert_utils.HexString2Bytes(value)[0]
	return
}

func parseMask(value string) bool {
	switch strings.ToLower(value) {
	case "true":
		return true
	default:
		return false
	}
}
