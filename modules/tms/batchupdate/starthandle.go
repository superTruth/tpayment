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
	logger := tlog.GetGoroutineLogger()

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 异步解析
	cCp := ctx.Copy()
	goroutine.Go(func() {
		tlog.SetGoroutineLogger(logger) // 切换协程，承接log
		StartUpdate(cCp, req.ID)
	})

	modules.BaseSuccess(ctx, nil)
}

func StartUpdate(ctx *gin.Context, id uint64) {
	logger := tlog.GetGoroutineLogger()

	// 获取批次记录
	updateRecord, err := tms.GetBatchUpdateRecordById(id)
	if err != nil {
		logger.Error("GetBatchUpdateRecordById error->", err.Error())
		return
	}

	logger.Info("start batch update->", updateRecord.ID)

	// 获取批次的升级app
	_, apps, err := tms.GetAppsInDevice(updateRecord.ID,
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
		devices, err := tms.GetBatchUpdateDevices(ctx, updateRecord, uint64(i*OnePageSize), OnePageSize)
		if err != nil {
			logger.Error("GetBatchUpdateDevices error->", err.Error())

			updateRecord.Status = "fail"
			updateRecord.UpdateFailMsg = "GetBatchUpdateDevices DB error"
			_ = models.UpdateBaseRecord(updateRecord)
			return
		}

		logger.Info("search devices -->", len(devices))

		for j := 0; j < len(devices); j++ { // 每一个设备
			deviceInBatch, err := tms.DeviceInBatchDao.GetByBatchIDDeviceID(updateRecord.ID, devices[j].ID)
			if err != nil {
				updateRecord.Status = "fail"
				updateRecord.UpdateFailMsg = "GetByBatchIDDeviceID DB error"
				_ = models.UpdateBaseRecord(updateRecord)
				return
			}
			if deviceInBatch != nil { // 已经设定过了，跳过
				continue
			}

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

			// 添加到更新列表里面
			err = tms.DeviceInBatchDao.Create(&tms.DeviceInBatchUpdate{
				BatchID:  updateRecord.ID,
				DeviceID: devices[j].ID,
				Status:   tms.BatchUpdateStatusPending,
			})
			if err != nil {
				logger.Errorf("DeviceInBatchDao.Create fail: %s", err.Error())
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
