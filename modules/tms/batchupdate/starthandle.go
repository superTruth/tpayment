package batchupdate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/modules/tms/appindevice"
	"tpayment/pkg/goroutine"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func StartHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(modules.BaseIDRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	// 异步解析
	goroutine.Go(func() {
		StartUpdate(ctx, req.ID)
	}, ctx)

	return nil
}

func StartUpdate(ctx echo.Context, id uint) {
	logger := tlog.GetLogger(ctx)

	// 获取批次记录
	updateRecord, err := tms.GetBatchUpdateRecordById(models.DB(), ctx, id)
	if err != nil {
		logger.Error("GetBatchUpdateRecordById error->", err.Error())
		return
	}

	// 获取批次的升级app
	_, apps, err := tms.GetAppsInDevice(models.DB(), ctx, updateRecord.ID,
		tms.AppInDeviceExternalIdTypeBatchUpdate, 0, 1000)
	if err != nil {
		logger.Error("GetAppsInDevice error->", err.Error())
		return
	}

	// 获取匹配的设备
	const OnePageSize = 1000
	for i := 0; ; i++ {
		devices, err := tms.GetBatchUpdateDevices(models.DB(), ctx, updateRecord, i*OnePageSize, (i+1)*OnePageSize)
		if err != nil {
			logger.Error("GetBatchUpdateDevices error->", err.Error())
			return
		}

		for j := 0; j < len(devices); j++ { // 每一个设备
			for k := 0; k < len(apps); k++ { // 每一个app
				ret := appindevice.SmartAddAppInDevice(ctx, devices[j], apps[k])
				if ret != conf.SUCCESS {
					logger.Error("config fail: ", ret.String(),
						", source file->", apps[k].AppFileId,
						", dest device->", devices[j].DeviceSn)
					return
				}
			}
		}
	}

}
