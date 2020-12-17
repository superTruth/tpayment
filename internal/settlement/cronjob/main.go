package main

import (
	"fmt"
	"time"
	"tpayment/conf"
	"tpayment/internal/basekey"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/payment/merchantaccount"
	"tpayment/pkg/tlog"

	"github.com/google/uuid"
)

const (
	MaxPage       = 1000
	MaxRetryTimes = 3
)

// 定时一个小时执行一次，把这个时间需要结算的收单里面的MID全部入消息队列
func main() {
	// 初始化日志
	requestId := uuid.New().String()
	logger := tlog.NewLog(requestId)
	defer logger.Destroy()

	logger.Info("start settlement cron job")

	conf.InitConfigData()

	models.InitDB()

	basekey.Init() // 初始化基础秘钥

	// 查找所有需要结算的收单
	hourStr := fmt.Sprintf("%02d", time.Now().Hour())
	logger.Info("match time->", hourStr)

	acqBean := &agency.Acquirer{
		BaseModel: models.BaseModel{
			Db: models.DB(),
		},
	}
	acquirers, err := acqBean.GetNeedSettlement(hourStr)
	if err != nil {
		logger.Error("GetNeedSettlement fail->", err.Error())
		return
	}

	logger.Info("match acquirer len->", len(acquirers))
	if len(acquirers) == 0 {
		logger.Info("no acquirer need settlement")
		return
	}

	// 遍历全部acquirer
	maBean := merchantaccount.MerchantAccount{
		BaseModel: models.BaseModel{
			Db: models.DB(),
		},
	}
	for _, acq := range acquirers {
		logger.Info("loop acquirer->", acq.Name)

		offset := uint64(0)
		retryTime := 0
		for {
			accounts, err := maBean.GetByAcquirerID(acq.ID, offset, MaxPage)
			if err != nil {
				retryTime++
				if retryTime > MaxRetryTimes {
					logger.Error("retry over time acq id->", acq.ID, ", offset", offset)
					return
				}
				time.Sleep(time.Second * 5)
			}
			retryTime = 0
			offset += MaxPage

			// TODO 全部入队
			for _, account := range accounts {
				logger.Info("push mq merchant->", account.ID, ", ", account.MID)
			}

			// 完成
			if len(accounts) < MaxPage {
				break
			}
		}
	}

}
