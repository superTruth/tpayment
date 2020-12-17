package common

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/bank_common"
	"tpayment/models/payment/acquirer"
	"tpayment/pkg/algorithmutils"
	"tpayment/pkg/iso8583"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/mix_utils"
)

func AddPlainHeader(req *acquirer_impl.SaleRequest, config *Config, msg *iso8583.Message) ([]byte, error) {
	headerBytes := convert_utils.HexString2Bytes(config.AcquirerConfig.TPDU + "600100180301")

	iso8583Bytes, err := msg.FormMessage()
	if err != nil {
		return nil, errors.New("FormMessage fail->" + err.Error())
	}

	lenBytes := convert_utils.Long2BytesHex(uint64(len(headerBytes)+len(iso8583Bytes)), 2)

	sb := bytes.Buffer{}

	sb.Write(lenBytes)
	sb.Write(headerBytes)
	sb.Write(iso8583Bytes)

	return sb.Bytes(), nil
}

func Encrypt(req *acquirer_impl.SaleRequest, config *Config, msg *iso8583.Message) ([]byte, error) {
	head1 := fmt.Sprintf("%s602100180321", config.AcquirerConfig.TPDU)
	head1Byte := convert_utils.HexString2Bytes(head1)
	head2 := fmt.Sprintf("%015s%08s",
		req.TxqReq.PaymentProcessRule.MerchantAccount.MID,
		req.TxqReq.PaymentProcessRule.MerchantAccount.Terminal.TID)

	iso8583Bytes, err := msg.FormMessage()
	if err != nil {
		return nil, errors.New("FormMessage fail->" + err.Error())
	}

	tak := FindKey(req, bank_common.TAK)
	tdk := FindKey(req, bank_common.TDK)
	if tak == nil || tdk == nil {
		return nil, errors.New("can't find key")
	}
	// 添加mac
	macBytes, err := algorithmutils.CalcCUPMac(
		mix_utils.BytesArrayCopyArrange(iso8583Bytes, 0,
			len(iso8583Bytes)-8), convert_utils.HexString2Bytes(tak.Value))
	if err != nil {
		return nil, errors.New("CalcCUPMac fail->" + err.Error())
	}
	mix_utils.BytesArrayCopy(macBytes, 0, iso8583Bytes, len(iso8583Bytes)-8, 8)

	// 添加明文数据长度
	plainLen := convert_utils.Long2BytesHex(uint64(len(iso8583Bytes)), 2)
	iso8583Bytes = mix_utils.MergeBytesArray(plainLen, iso8583Bytes)

	// 加密一下报文
	iso8583Bytes, _ = algorithmutils.EncryptDesECB(iso8583Bytes, convert_utils.HexString2Bytes(tdk.Value))

	lenBytes := convert_utils.Long2BytesHex(uint64(len(head1Byte)+len(head2)+len(iso8583Bytes)), 2)

	sb := bytes.Buffer{}

	sb.Write(lenBytes)
	sb.Write(head1Byte)
	sb.Write([]byte(head2))
	sb.Write(iso8583Bytes)

	return sb.Bytes(), nil
}

func Decrypt(data []byte, req *acquirer_impl.SaleRequest) ([]byte, error) {
	tak := FindKey(req, bank_common.TAK)
	tdk := FindKey(req, bank_common.TDK)
	if tak == nil || tdk == nil {
		return nil, errors.New("can't find key")
	}

	ret, _ := algorithmutils.DecryptDesECB(data, convert_utils.HexString2Bytes(tdk.Value))

	if len(ret) < 5 {
		return nil, errors.New("receive plain data too short:" + strconv.Itoa(len(ret)))
	}

	plainLen := int(convert_utils.BytesHex2Long(ret, 0, 2))

	if (plainLen + 2) > len(ret) {
		return nil, errors.New("receive plain data too short 2:" + strconv.Itoa(len(ret)))
	}

	ret = mix_utils.BytesArrayCopyArrange(ret, 2, 2+plainLen)

	return ret, nil
}

func FindKey(req *acquirer_impl.SaleRequest, keyType string) *acquirer.Key {
	if len(req.Keys) == 0 {
		return nil
	}

	for i := 0; i < len(req.Keys); i++ {
		if req.Keys[i].Type == keyType {
			return req.Keys[i]
		}
	}

	return nil
}
