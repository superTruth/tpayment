package logon

import (
	"errors"
	"fmt"
	"strings"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/bank_common"
	"tpayment/internal/bank_service/sic/internal/common"
	"tpayment/models/payment/acquirer"
	"tpayment/pkg/algorithmutils"
	"tpayment/pkg/emv/tlv"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils/convert_utils"
)

type DownloadTMK struct {
	Req         *acquirer_impl.SaleRequest
	Config      *common.Config
	ExchangeKey string
	PukModulus  string
	PukExponent string
	TMKEn       string
}

func (s *DownloadTMK) Handle() (*acquirer_impl.SaleResponse, conf.ResultCode) {
	logger := tlog.GetGoroutineLogger()

	s.ExchangeKey = generateKekPlainHexStr()

	// 打包数据
	sendBytes, err := s.packageMsg()
	if err != nil {
		logger.Error("packageMsg fail->", err.Error())
		return nil, conf.ParameterError
	}

	// 数据交换
	retBytes, errCode := common.Exchange(s.Req, s.Config, sendBytes)
	if errCode != conf.Success {
		logger.Warn("exchange fail->", errCode.String())
		return nil, conf.CantReachAcquirer
	}

	// 解包数据
	resp, err := s.parseMsg(retBytes)
	if err != nil {
		logger.Warn("parseMsg fail->", err.Error())
		return nil, conf.UnknownError
	}

	if resp.TxnResp.CreditCardBean.ResponseCode != "00" {
		logger.Warn("reject by acquirer, response code is ->", resp.TxnResp.CreditCardBean.ResponseCode)
		return resp, conf.RejectByAcquirer
	}

	if err = s.handle(resp); err != nil {
		logger.Warn("decode fail->", err.Error())
		return resp, conf.RejectByAcquirer
	}

	return resp, conf.Success
}

func (s *DownloadTMK) handle(resp *acquirer_impl.SaleResponse) error {
	tmk, desErr := algorithmutils.DecryptDesECB(convert_utils.HexString2Bytes(s.TMKEn), convert_utils.HexString2Bytes(s.ExchangeKey))
	if desErr != nil {
		return desErr
	}

	resp.Keys = append(resp.Keys, &acquirer.Key{
		Type:  bank_common.TMK,
		Value: convert_utils.Bytes2HexString(tmk),
	})

	return nil
}

func (s *DownloadTMK) packageMsg() ([]byte, error) {
	logger := tlog.GetGoroutineLogger()

	msg, err := common.SetCommonDataToMsg(s.Req, s.Config,
		[]byte{11, 41, 42, 60, 62})
	if err != nil {
		return nil, errors.New("parameter error->" + err.Error())
	}

	_ = msg.SetMessageType("0800")

	// field60
	field60Sb := strings.Builder{}
	field60Sb.WriteString("99")
	field60Sb.WriteString(fmt.Sprintf("%06d", s.Req.TxqReq.CreditCardBean.BatchNum))
	field60Sb.WriteString("003")
	_ = msg.SetFieldValue(60, field60Sb.String())

	// field62
	tmkKeyFormat := generateKekPKCS1HexStr(s.ExchangeKey)
	keyTmp := convert_utils.HexString2Bytes(tmkKeyFormat)

	tmkEnByPk, err := algorithmutils.RsaPublicEncryptionX(keyTmp, s.PukModulus, s.PukExponent)
	if err != nil {
		return nil, errors.New("RsaPublicEncryptionX fail->" + err.Error())
	}
	tmkEnByPkStr := convert_utils.Bytes2HexString(tmkEnByPk)

	mapData := make(map[string]string) // format to tlv data
	mapData["DF99"] = tmkEnByPkStr
	mapData["9F06"] = "DF00000003" // Landi key
	mapData["9F22"] = "01"         // Landi key

	_ = msg.SetFieldValue(62, tlv.FormatFromMap(mapData))

	// field63
	field63Sb := strings.Builder{}
	field63Sb.WriteString("001")
	field63Sb.WriteString(s.Req.TxqReq.DeviceID)
	_ = msg.SetFieldValue(63, field63Sb.String())

	_ = msg.SetFieldValue(62, "9F0605DF000000039F220101")

	logger.Info("package msg->", msg.String())

	sendBytes, err := common.AddPlainHeader(s.Req, s.Config, msg)
	if err != nil {
		return nil, errors.New("encrypt error->" + err.Error())
	}
	return sendBytes, nil
}

func (s *DownloadTMK) parseMsg(msgData []byte) (*acquirer_impl.SaleResponse, error) {
	resp, msg, err := common.GetCommonDataFromMsg(msgData, 0)
	if err != nil {
		return nil, errors.New("getCommonDataFromMsg fail->" + err.Error())
	}

	// TODO 测试数据 38A216013DC29D734F64FB7FFD624507
	tmkEn, err := algorithmutils.EncryptDesECB(
		convert_utils.HexString2Bytes("38A216013DC29D734F64FB7FFD624507"),
		convert_utils.HexString2Bytes(s.ExchangeKey))
	if err != nil {
		return nil, err
	}
	resp.TxnResp.CreditCardBean.ResponseCode = "00"
	s.TMKEn = convert_utils.Bytes2HexString(tmkEn)
	//resp.Keys = append(resp.Keys, &acquirer.Key{
	//	Type:  bank_common.TMK,
	//	Value: "38A216013DC29D734F64FB7FFD624507",
	//})
	return resp, nil

	// field 39
	field39, _ := msg.GetFieldValue(39)
	if field39 == "" {
		return nil, errors.New("f39 is empty")
	}

	if field39 != "00" {
		resp.TxnResp.CreditCardBean.ResponseCode = field39
		return resp, nil
	}

	// field62
	field62, _ := msg.GetFieldValue(62)
	if len(field62) < 10 {
		return nil, errors.New("f62 format error")
	}

	field62 = field62[2:] // 去掉第一个字节

	tlvMap, err := tlv.Parse2Map(field62, true)
	if err != nil {
		return nil, errors.New("f62 parse to tlv fail->" + err.Error())
	}

	tmk, ok := tlvMap["DF02"]
	if !ok {
		return nil, errors.New("can't find tmk")
	}

	s.TMKEn = tmk
	//resp.Keys = append(resp.Keys, &acquirer.Key{
	//	Type:  bank_common.TMK,
	//	Value: tmk,
	//})

	return resp, nil
}
