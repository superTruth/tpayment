package qrcode

import (
	"encoding/base64"
	"errors"
	"strings"
	"tpayment/pkg/emv/tlv"
	"tpayment/pkg/paymentmethod/decodecardnum/creditcard"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/format_utils"
)

type EMVQRDecodeContent struct {
	CardNum string
	Track2  string
	CardSn  string
	ICCData string
}

func DecodeEmvQR(qrCode string) (*EMVQRDecodeContent, error) {
	// 部分银联二维码不是标准的base64格式，需要填充
	qrCode = format_utils.AppendString(qrCode, (len(qrCode)+3)/4*4, false, '=')
	qrCodeByte, err := base64.StdEncoding.DecodeString(qrCode)
	if err != nil {
		return nil, errors.New("qr not base64 format")
	}

	dataMap, err := tlv.Parse2Map(convert_utils.Bytes2HexString(qrCodeByte), false)
	if err != nil {
		return nil, errors.New("qr not tlv format")
	}

	// 提取tag 61
	tag61, ok := dataMap["61"]
	if !ok {
		tag61, ok = dataMap["65"]
		if !ok {
			return nil, errors.New("can't find tag61")
		}
	}

	ret := new(EMVQRDecodeContent)
	// 再次解析tag 61
	dataMap, err = tlv.Parse2Map(tag61, false)
	if err != nil {
		return nil, errors.New("qr not tlv format")
	}

	ret.CardNum = dataMap["5A"]
	ret.Track2 = dataMap["57"]

	if ret.CardNum == "" {
		if ret.Track2 == "" {
			return nil, errors.New("can't find card number")
		}
		ret.CardNum, err = creditcard.GetCardNumFromTrack2(ret.Track2)
		if err != nil {
			return nil, errors.New("can't find card number from tk2")
		}
	}

	// 去除非数字
	ret.CardNum = strings.TrimFunc(ret.CardNum, func(r rune) bool {
		return r < '0' || r > '9'
	})

	// card SN

	// TODO

	return ret, nil
}
