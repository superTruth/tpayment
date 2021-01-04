package common

import (
	"time"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/bank_service/bank_common/communicate/socket"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils/convert_utils"
)

func Exchange(req *acquirer_impl.SaleRequest, config *Config, sendBytes []byte) ([]byte, conf.ResultCode) {
	logger := tlog.GetGoroutineLogger()
	con := socket.Generate(&socket.InitObject{
		URL:      config.AcquirerConfig.URL,
		SSLModel: socket.SSLEnableSkip,
	})
	defer con.Disconnect(time.Second * 1) // 使用完，需要断开连接

	logger.Info("start ExchangeMsg....->", convert_utils.Bytes2HexString(sendBytes))
	retBytes, errCode := con.ExchangeMsg(sendBytes, true, &Cn8583Protocol{},
		time.Second*50)
	logger.Info("end ExchangeMsg->", errCode, ", body->", convert_utils.Bytes2HexString(retBytes))
	switch errCode {
	case socket.Success:
		return retBytes, conf.Success
	case socket.ConnectFail:
		return nil, conf.CantReachAcquirer
	default:
		return nil, conf.Reversal
	}
}
