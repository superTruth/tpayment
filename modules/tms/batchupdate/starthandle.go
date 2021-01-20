package batchupdate

import (
	"fmt"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/modules/tms/appindevice"
	"tpayment/pkg/goroutine"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func StartHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	modules.BaseSuccess(ctx, nil)

	// 异步解析
	goroutine.Go(func() {
		StartUpdate(ctx, req.ID)
	}, ctx)
}

func StartUpdate(ctx *gin.Context, id uint64) {
	logger := tlog.GetLogger(ctx)

	// 获取批次记录
	updateRecord, err := tms.GetBatchUpdateRecordById(models.DB(), ctx, id)
	if err != nil {
		logger.Error("GetBatchUpdateRecordById error->", err.Error())
		return
	}

	logger.Info("start batch update->", updateRecord.ID)

	// 获取批次的升级app
	_, apps, err := tms.GetAppsInDevice(models.DB(), ctx, updateRecord.ID,
		tms.AppInDeviceExternalIdTypeBatchUpdate, 0, 1000)
	if err != nil {
		logger.Error("GetAppsInDevice error->", err.Error())
		return
	}

	logger.Info("GetAppsInDevice len->", len(apps))

	updateRecord.Status = "updating"
	_ = models.UpdateBaseRecord(updateRecord)

	// 获取匹配的设备
	const OnePageSize = 1000
	for i := 0; ; i++ {
		devices, err := tms.GetBatchUpdateDevices(models.DB(), ctx, updateRecord, uint64(i*OnePageSize), uint64((i+1)*OnePageSize))
		if err != nil {
			logger.Error("GetBatchUpdateDevices error->", err.Error())

			updateRecord.Status = "fail"
			updateRecord.UpdateFailMsg = "GetBatchUpdateDevices DB error"
			_ = models.UpdateBaseRecord(updateRecord)
			return
		}

		logger.Info("search devices -->", len(devices))

		for j := 0; j < len(devices); j++ { // 每一个设备
			logger.Info("device updating -->", devices[j].DeviceSn, ",", devices[j].ID)
			for k := 0; k < len(apps); k++ { // 每一个app
				ret := appindevice.SmartAddAppInDevice(ctx, devices[j], apps[k])
				if ret != conf.Success {
					errStr := fmt.Sprint("config fail: ", ret.String(),
						", source file->", apps[k].AppFileId,
						", dest device->", devices[j].DeviceSn)
					logger.Error(errStr)

					updateRecord.Status = "fail"
					updateRecord.UpdateFailMsg = "dest device->" + devices[j].DeviceSn
					_ = models.UpdateBaseRecord(updateRecord)
					return
				}
			}
		}

		if len(devices) < OnePageSize {
			break
		}
	}

	updateRecord.Status = "Done"
	updateRecord.UpdateFailMsg = " "
	_ = models.UpdateBaseRecord(updateRecord)
}
