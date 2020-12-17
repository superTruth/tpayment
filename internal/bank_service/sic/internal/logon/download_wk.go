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
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/mix_utils"
)

type DownloadWK struct {
	Req    *acquirer_impl.SaleRequest
	Config *common.Config
}

func (s *DownloadWK) Handle() (*acquirer_impl.SaleResponse, conf.ResultCode) {
	logger := tlog.GetGoroutineLogger()

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
		return nil, errCode
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

	// 解密秘钥
	tmk := common.FindKey(s.Req, bank_common.TMK)
	resp.Keys, err = decryptWK(tmk, resp.Keys)
	if err != nil {
		logger.Warn("decrypt key error->", err)
		return nil, conf.RejectByAcquirer
	}

	return resp, conf.Success
}

func decryptWK(tmk *acquirer.Key, wks []*acquirer.Key) ([]*acquirer.Key, error) {
	tmkByte := convert_utils.HexString2Bytes(tmk.Value)
	for i := range wks {
		orgBytes := convert_utils.HexString2Bytes(wks[i].Value)
		keyPlain, err := algorithmutils.DecryptDesECB(
			mix_utils.BytesArrayCopyArrange(orgBytes, 0, len(orgBytes)-4), tmkByte)
		if err != nil {
			return nil, errors.New("DecryptDesECB fail->" + err.Error())
		}
		kcv := mix_utils.BytesArrayCopyArrange(orgBytes, len(orgBytes)-4, len(orgBytes))
		if !algorithmutils.CheckKCV(keyPlain, kcv) {
			return nil, errors.New("kcv error")
		}
		wks[i].Value = convert_utils.Bytes2HexString(keyPlain)
	}
	return wks, nil
}

func (s *DownloadWK) packageMsg() ([]byte, error) {
	logger := tlog.GetGoroutineLogger()

	msg, err := common.SetCommonDataToMsg(s.Req, s.Config,
		[]byte{11, 41, 42, 60, 63})
	if err != nil {
		return nil, errors.New("parameter error->" + err.Error())
	}

	_ = msg.SetMessageType("0800")

	// field60
	field60Sb := strings.Builder{}
	field60Sb.WriteString("00")
	field60Sb.WriteString(fmt.Sprintf("%06d", s.Req.TxqReq.CreditCardBean.BatchNum))
	field60Sb.WriteString("003")
	_ = msg.SetFieldValue(60, field60Sb.String())

	// field63
	_ = msg.SetFieldValue(63, "001") // CashierNo

	logger.Info("package msg->", msg.String())

	sendBytes, err := common.AddPlainHeader(s.Req, s.Config, msg)
	if err != nil {
		return nil, errors.New("encrypt error->" + err.Error())
	}
	return sendBytes, nil
}

func (s *DownloadWK) parseMsg(msgData []byte) (*acquirer_impl.SaleResponse, error) {
	resp, msg, err := common.GetCommonDataFromMsg(msgData, 0)
	if err != nil {
		return nil, errors.New("getCommonDataFromMsg fail->" + err.Error())
	}

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
	field60, _ := msg.GetFieldValue(60)
	field62, _ := msg.GetFieldValue(62)
	if len(field60) != 11 {
		return nil, errors.New("f60 format error")
	}
	//batchNum := field60[2:8]  // batch num
	keyType := field60[8:]

	if len(field62) != 122 {
		return nil, errors.New("f62 format error")
	}
	field62 = field62[2:] // 去掉第一个字节

	var (
		tpk, tak, tdk string
	)
	switch keyType {
	case "001":
		if len(field62) != 48 {
			return nil, errors.New("key type 001, but key len not 24")
		}
		tpk = field62[:24]
		tak = field62[24:48]
	case "003":
		if len(field62) != 120 {
			return nil, errors.New("key type 003, but key len not 40")
		}
		tpk = field62[:40]
		tak = field62[40:80]
		tdk = field62[80:120]
		break
	default:
		return nil, errors.New("unknow key type:" + keyType)
	}

	resp.Keys = append(resp.Keys, &acquirer.Key{
		Type:  bank_common.TPK,
		Value: tpk,
	})
	resp.Keys = append(resp.Keys, &acquirer.Key{
		Type:  bank_common.TAK,
		Value: tak,
	})
	resp.Keys = append(resp.Keys, &acquirer.Key{
		Type:  bank_common.TDK,
		Value: tdk,
	})

	return resp, nil
}
