package tlv

import (
	"errors"
	"strconv"
	"tpayment/pkg/utils/convert_utils"
)

type Node struct {
	Tag   string
	Value string
}

func Parse(dataStr string, isPBOC bool) ([]*Node, error) {
	var dataList []*Node

	dataByte := convert_utils.HexString2Bytes(dataStr)
	if len(dataByte) == 0 {
		return dataList, errors.New("source data is empty")
	}

	for index := 0; index < len(dataByte); {
		node, offset, err := parseOneNode(dataByte, index, isPBOC)
		if err != nil {
			return dataList, err
		}

		index += offset
		dataList = append(dataList, node)
	}

	return dataList, nil
}

func Parse2Map(dataStr string, isPBOC bool) (map[string]string, error) {
	nodes, err := Parse(dataStr, isPBOC)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)

	for _, node := range nodes {
		if _, ok := ret[node.Tag]; ok {
			continue
		}
		ret[node.Tag] = node.Value
	}

	return ret, nil
}

func parseOneNode(orgData []byte, offset int, isPBOC bool) (*Node, int, error) {
	orgLen := len(orgData)
	index := offset
	// Get Tag length
	tagLen := 1
	if (index+1 < orgLen) && (0x1F == (orgData[index] & 0x1F)) {
		if !isPBOC && 0x80 == (orgData[index+1]&0x80) {
			tagLen = 3
		} else {
			tagLen = 2
		}
	} else {
		tagLen = 1
	}

	// 提取Tag
	if index+tagLen > orgLen {
		return nil, 0, errors.New("tag len out of range->" + strconv.Itoa(orgLen) + strconv.Itoa(index+tagLen))
	}
	tag := orgData[index : index+tagLen]

	index += tagLen
	// 提取value Len
	if index > orgLen {
		return nil, 0, errors.New("can't find value len len->" + strconv.Itoa(orgLen) + strconv.Itoa(index))
	}

	valueLenLen := 1
	if (orgData[index] & 0x80) == 0x80 { // 多数据长度
		valueLenLen = int(orgData[index] & 0x7F)
		index += 1 // 跳过长度位
	}

	if index+valueLenLen > orgLen {
		return nil, 0, errors.New("can't find value len->" + strconv.Itoa(orgLen) + strconv.Itoa(index+valueLenLen))
	}
	valueLen := int(convert_utils.BytesHex2Long(orgData, index, valueLenLen))
	index += valueLenLen

	// 提取value
	if (index + valueLen) > orgLen {
		return nil, 0, errors.New("can't find value->" + strconv.Itoa(orgLen) + strconv.Itoa(index))
	}
	valueByte := orgData[index : index+valueLen]
	index += valueLen

	return &Node{
		Tag:   convert_utils.Bytes2HexString(tag),
		Value: convert_utils.Bytes2HexString(valueByte),
	}, index - offset, nil
}
