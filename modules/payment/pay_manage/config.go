package pay_manage

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetPaymentConfig(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseQueryRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 获取device id
	device, err := tms.DeviceInfoDao.GetBySn(req.DeviceSN)
	if err != nil {
		logger.Error("GetBySn sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if device == nil {
		logger.Info("record not found->", req.DeviceSN)
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 查询merchant device
	merchantDevice, err := merchant.DeviceInMerchantDao.GetByMerchantIdAndDeviceID(req.MerchantId, device.ID)
	if err != nil {
		logger.Error("GetByMerchantIdAndDeviceID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if merchantDevice == nil {
		logger.Info("record not found -> merchant id:", req.MerchantId, " device ID->", device.ID)
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	if req.Limit > conf.MaxQueryCount || req.Limit == 0 { // 一次性不能搜索太多数据
		req.Limit = conf.MaxQueryCount
	}

	// 查数据
	total, dataRet, err := merchant.QueryPaymentSettingInDeviceRecord(models.DB(), ctx, merchantDevice.ID, req.Offset, req.Limit, req.Filters)
	if err != nil {
		logger.Error("QueryBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	// 查找所有的acquirer config
	for i := 0; i < len(dataRet); i++ {
		dataRet[i].AcquirerConfig, err = agency.AcquirerDao.Get(dataRet[i].AcquirerId)
		if err != nil {
			logger.Error("agency.AcquirerDao.Get sql error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return
		}
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)
}
