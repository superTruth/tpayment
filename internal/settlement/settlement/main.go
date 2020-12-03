package main

import (
	"time"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/internal/acquirer_impl/factory"
	"tpayment/internal/basekey"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
	"tpayment/models/payment/record"
	"tpayment/pkg/tlog"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// 消费结算
func main() {

	conf.InitConfigData()

	models.InitDB()

	basekey.Init() // 初始化基础秘钥
}

const (
	maxSettlementTime = time.Minute * 5
	maxRetryTime      = 10
)

func settlement(maID uint) error {
	var (
		err       error
		errorCode conf.ResultCode
	)

	requestId := uuid.New().String()
	logger := tlog.NewLog(requestId)
	defer logger.Destroy()

	// 获取acquirer
	merchantBean := &merchantaccount.MerchantAccount{
		BaseModel: models.BaseModel{
			Db: models.DB(),
		},
	}
	merchantBean, err = merchantBean.Get(maID)
	if err != nil {
		logger.Error("merchantBean.Get err->", maID)
		return err
	}
	if merchantBean == nil {
		logger.Error("merchantBean.Get can't get record->", maID)
		return nil
	}

	acqBean := &agency.Acquirer{
		BaseModel: models.BaseModel{
			Db: models.DB(),
		},
	}
	acqBean, err = acqBean.Get(merchantBean.AcquirerID)
	if err != nil {
		logger.Error("acqBean.Get err->", maID)
		return err
	}
	if acqBean == nil {
		logger.Error("acqBean.Get can't get record->", maID)
		return nil
	}

	// 匹配acquirer实现
	acquirerImpl, ok := factory.AcquirerImpls[acqBean.ImplName]
	if !ok {
		logger.Warn("can't find acquirer impl->", acqBean.Name)
		return nil
	}
	settlementInMID, ok := acquirerImpl.(acquirer_impl.ISettlementInMID)
	if ok {
		// MID维度
		logger.Info("acquirer ", acqBean.Name, " settlement by mid")
		for i := 0; i < maxRetryTime; i++ {
			errorCode = settlementInMID.SettlementInMID(acqBean, merchantBean,
				func(f ...func(*gorm.DB) *gorm.DB) *gorm.DB {
					t := time.Now()
					// 添加一层限制，防止出现更新别的数据情况
					return models.DB().Model(&record.TxnRecord{}).
						Scopes(f...).
						Where("merchant_account_id=?", maID).
						Update("acquirer_settlement_at", &t)
				})
			if errorCode == conf.Success {
				break
			}
			time.Sleep(time.Second * 10)
		}

		return nil
	}

	// TID维度
	settlementInTID, ok := acquirerImpl.(acquirer_impl.ISettlementInTID)
	if ok {
		// TID维度
		logger.Info("acquirer ", acqBean.Name, " settlement by tid")

		// 查询TID
		tid := &acquirer.Terminal{
			BaseModel: models.BaseModel{
				Db: models.DB(),
			},
		}
		tids, err := tid.GetByMID(maID)
		if err != nil {
			logger.Error("tid.GetByMID fail->", maID)
			return err
		}

		if len(tids) == 0 {
			logger.Warn("merchant ", maID, " not config tid")
			return nil
		}

		retryTime := 0
		// 遍历结算所有TID，进行结算
		for i := 0; ; {
			if tids[i] == nil { // 已经结算完成的，不要再次出现
				goto NextOne
			}

			// 锁定TID
			tids[i].Db = models.DB()
			errorCode = tids[i].Lock(maxSettlementTime)
			if errorCode != conf.Success { // 锁定失败
				goto NextOne
			}

			// 开始结算
			errorCode = settlementInTID.SettlementInTID(acqBean, merchantBean, tids[i],
				func(f ...func(*gorm.DB) *gorm.DB) *gorm.DB {
					t := time.Now()
					// 添加一层限制，防止出现更新别的数据情况
					return models.DB().Model(&record.TxnRecord{}).
						Scopes(f...).
						Where("merchant_account_id=? AND terminal_id=?", maID, tids[i].ID).
						Update("acquirer_settlement_at", &t)
				})

			// 解锁TID
			_ = tids[i].UnLock()

			if errorCode == conf.Success {
				tids[i] = nil // 成功的话，就删除掉这个数据
				goto NextOne
			}

		NextOne:
			i = (i + 1) % len(tids) // 循环
			if i == 0 {             //完成一个循环
				if !checkTidExit(tids) { // 正常完成所有记录
					logger.Info("complete settlement success->", maID)
					return nil
				}
				retryTime++
				if retryTime > maxRetryTime { // 超过最大尝试循环次数还是无法成功
					logger.Error("can't settlement success by tid for mid ->", maID)
					return nil
				}

				time.Sleep(time.Second * 10) // 10秒后再次尝试
			}
		}

	}

	logger.Error("the acquirer not support settlement->", acqBean.Name)
	return nil
}

func checkTidExit(tids []*acquirer.Terminal) bool {
	for _, tid := range tids {
		if tid != nil {
			return true
		}
	}
	return false
}
