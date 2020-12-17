package logon

import (
	"errors"
	"fmt"
	"strings"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/sic/internal/common"
	"tpayment/pkg/emv/tlv"
	"tpayment/pkg/tlog"
)

type DownloadPuk struct {
	Req         *acquirer_impl.SaleRequest
	Config      *common.Config
	PukModulus  string
	PukExponent string
}

func (s *DownloadPuk) Handle() (*acquirer_impl.SaleResponse, conf.ResultCode) {
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
		return nil, conf.CantReachAcquirer
	}

	// 解包数据
	resp, err := s.parseMsg(retBytes)
	if err != nil {
		logger.Error("packageMsg fail->", err.Error())
		return nil, conf.UnknownError
	}
	if resp.TxnResp.CreditCardBean.ResponseCode != "00" {
		logger.Warn("reject by acquirer, response code is ->", resp.TxnResp.CreditCardBean.ResponseCode)
		return resp, conf.RejectByAcquirer
	}

	return resp, conf.Success
}

func (s *DownloadPuk) packageMsg() ([]byte, error) {
	logger := tlog.GetGoroutineLogger()

	msg, err := common.SetCommonDataToMsg(s.Req, s.Config,
		[]byte{41, 42, 60, 62})
	if err != nil {
		return nil, errors.New("parameter error->" + err.Error())
	}

	_ = msg.SetMessageType("0800")

	// field60
	field60Sb := strings.Builder{}
	field60Sb.WriteString("96")
	field60Sb.WriteString(fmt.Sprintf("%06d", s.Req.TxqReq.CreditCardBean.BatchNum))
	field60Sb.WriteString("400")
	_ = msg.SetFieldValue(60, field60Sb.String())

	_ = msg.SetFieldValue(62, "9F0605DF000000039F220101")

	logger.Info("package msg->", msg.String())

	sendBytes, err := common.AddPlainHeader(s.Req, s.Config, msg)
	if err != nil {
		return nil, errors.New("encrypt error->" + err.Error())
	}
	return sendBytes, nil
}

func (s *DownloadPuk) parseMsg(msgData []byte) (*acquirer_impl.SaleResponse, error) {
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
	field62, _ := msg.GetFieldValue(62)
	if len(field62) < 10 {
		return nil, errors.New("f62 format error")
	}

	field62 = field62[2:] // 去掉第一个字节

	tlvMap, err := tlv.Parse2Map(field62, true)
	if err != nil {
		return nil, errors.New("f62 parse to tlv fail->" + err.Error())
	}

	ok := false
	s.PukModulus, ok = tlvMap["DF02"]
	if !ok {
		return nil, errors.New("can't find modulus")
	}
	s.PukExponent, ok = tlvMap["DF04"]
	if !ok {
		return nil, errors.New("can't find exponent")
	}

	return resp, nil
}
