package appindevice

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	tms2 "tpayment/modules/tms"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func QueryHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseQueryRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	if req.Limit > conf.MaxQueryCount { // 一次性不能搜索太多数据
		req.Limit = conf.MaxQueryCount
	}

	// 权限判断
	deviceBean, err := tms.GetDeviceByID(models.DB(), ctx, req.DeviceId)
	if err != nil {
		logger.Info("GetDeviceByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if err := tms2.CheckPermission(ctx, deviceBean); err != nil {
		logger.Error("CheckPermission fail->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	total, dataRet, err := tms.GetAppsInDevice(models.DB(), ctx, req.DeviceId,
		tms.AppInDeviceExternalIdTypeDevice, req.Offset, req.Limit)
	if err != nil {
		logger.Info("QueryAppInDeviceRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	// 整理数据
	for i := 0; i < len(dataRet); i++ {
		// 等待安装或者等待卸载时， 需要把配置的app信息显示替换掉已存在的显示出来
		if (dataRet[i].Status == conf.TmsStatusPendingInstall) ||
			(dataRet[i].Status == conf.TmsStatusPendingUninstalled) {
			if dataRet[i].App != nil {
				dataRet[i].PackageId = dataRet[i].App.PackageId
				dataRet[i].Name = dataRet[i].App.Name
			}

			if dataRet[i].AppFile != nil {
				dataRet[i].VersionName = dataRet[i].AppFile.VersionName
				dataRet[i].VersionCode = dataRet[i].AppFile.VersionCode
			}
		}
	}

	ret := &modules.BaseQueryResponse{
		Total: total,
		Data:  dataRet,
	}

	modules.BaseSuccess(ctx, ret)
}
