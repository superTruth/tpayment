package standard

import (
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
	"tpayment/models/payment/record"
	"tpayment/pkg/tlog"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func (api *API) SettlementInTID(acq *agency.Acquirer, mid *merchantaccount.MerchantAccount,
	tid *acquirer.Terminal, dbFunc func(...func(*gorm.DB) *gorm.DB) *gorm.DB) conf.ResultCode {
	var (
		err error
	)
	// 初始化日志
	requestId := uuid.New().String()
	logger := tlog.NewLog(requestId)
	defer logger.Destroy()

	// 获取统计信息
	total, err := record.TxnRecordDao.GetSettlementTotal(mid.ID, tid.ID, tid.BatchNum)
	if err != nil {
		logger.Error("recordBean.GetSettlementTotal fail->", err)
		return conf.DBError
	}
	if len(total) == 0 {
		logger.Info("no record need settlement")
		return conf.Success
	}

	// TODO 发起结算

	// TODO 第一次批上送

	// TODO 中间批上送

	// TODO 最后一次批上送

	// TODO 批上送结束

	return conf.Success
}
