package appinbatchupdate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func QueryHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseQueryRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	if req.Limit > conf.MaxQueryCount { // 一次性不能搜索太多数据
		req.Limit = conf.MaxQueryCount
	}

	total, dataRet, err := tms.GetAppsInDevice(models.DB(), ctx, req.BatchId, tms.AppInDeviceExternalIdTypeBatchUpdate, req.Offset, req.Limit)
	if err != nil {
		logger.Info("QueryAppInDeviceRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	// 整理数据
	for i := 0; i < len(dataRet); i++ {
		// 等待安装或者等待卸载时， 需要把配置的app信息显示替换掉已存在的显示出来
		if (dataRet[i].Status == conf.TmsStatusPendingInstall) ||
			(dataRet[i].Status == conf.TmsStatusPendingUninstalled) {
			if dataRet[i].App != nil {
				dataRet[i].PackageId = dataRet[i].App.PackageId
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

	return nil
}
