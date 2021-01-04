package sale

import (
	"errors"
	"fmt"
	"strings"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/boc/internal/common"
	"tpayment/models/payment/record"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/tlog"
)

type Sale struct {
	Req    *acquirer_impl.SaleRequest
	Config *common.Config
}

func (s *Sale) Handle() (*acquirer_impl.SaleResponse, conf.ResultCode) {
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
		logger.Error("parseMsg fail->", err.Error())
		return nil, conf.Reversal
	}

	s.handle(resp)

	return resp, conf.Success
}

func (s *Sale) handle(resp *acquirer_impl.SaleResponse) {
	logger := tlog.GetGoroutineLogger()

	if resp.TxnResp.CreditCardBean.ResponseCode != "00" {
		logger.Warn("reject by acquirer, response code ->", resp.TxnResp.CreditCardBean.ResponseCode)
		resp.TxnResp.TransactionState = record.Fail
		resp.TxnResp.ErrorCode = string(conf.RejectByAcquirer)
		resp.TxnResp.ErrorDesc = conf.RejectByAcquirer.String()
		return
	}
	resp.AcquirerReconID = resp.TxnResp.AcquirerRRN
	logger.Info("txn success auth code:", resp.TxnResp.CreditCardBean.AuthCode, ", rrn->", resp.TxnResp.AcquirerRRN)
	resp.TxnResp.TransactionState = record.Success
}

func (s *Sale) packageMsg() ([]byte, error) {
	logger := tlog.GetGoroutineLogger()
	msg, err := common.SetCommonDataToMsg(s.Req, s.Config,
		[]byte{2, 3, 4, 11, 14, 22, 23, 25, 26, 35, 41, 42, 49, 52, 53, 55, 64})
	if err != nil {
		return nil, errors.New("parameter error->" + err.Error())
	}

	_ = msg.SetMessageType("0200")
	_ = msg.SetFieldValue(3, "190000")
	_ = msg.SetFieldValue(25, "82")
	_ = msg.SetFieldValue(59, "GS0610004161001050    0240113.327492          23.118519           SN008A8000002ON0201")

	// field60
	field60Sb := strings.Builder{}
	field60Sb.WriteString("22")
	field60Sb.WriteString(fmt.Sprintf("%06d", s.Req.TxqReq.CreditCardBean.BatchNum))
	field60Sb.WriteString("000")
	field60Sb.WriteString("6")
	field60Sb.WriteString("0")
	if s.Req.TxqReq.RealEntryType == conf.ContactLess {
		field60Sb.WriteString("1")
	} else {
		field60Sb.WriteString("0")
	}
	_ = msg.SetFieldValue(60, field60Sb.String())

	logger.Info("package msg->", msg.String())

	sendBytes, err := msg.FormMessageX(2, iso8583Define.Hex)

	if err != nil {
		return nil, errors.New("encrypt error->" + err.Error())
	}
	return sendBytes, nil
}

func (s *Sale) parseMsg(msgData []byte) (*acquirer_impl.SaleResponse, error) {
	resp, _, err := common.GetCommonDataFromMsg(msgData, 0)
	if err != nil {
		return nil, errors.New("getCommonDataFromMsg fail->" + err.Error())
	}

	return resp, nil
}
