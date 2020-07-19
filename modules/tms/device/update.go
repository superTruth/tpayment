package device

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

func UpdateHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(tms.DeviceInfo)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询是否已经存在的账号
	bean, err := tms.GetDeviceByID(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetDeviceByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if bean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	// 生成新账号
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	// 合并tags
	errorCode := mergeTags(ctx, req)
	if errorCode != conf.SUCCESS {
		modules.BaseError(ctx, errorCode)
		return errors.New(errorCode.String())
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}

// 合并tags
func mergeTags(ctx echo.Context, device *tms.DeviceInfo) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	dbTags, err := tms.QueryTags(models.DB(), ctx, device)

	if err != nil {
		logger.Warn("QueryTags fail->", err.Error())
		return conf.DBError
	}

	// 先找到需要删除的tag(在数据库中有，但是update没有的)
	logger.Info("find need delete record-> db len:", len(dbTags), ", tags len:", len(device.Tags))
	for i := 0; i < len(dbTags); i++ {
		findFlag := false

		for j := 0; j < len(device.Tags); j++ {
			logger.Info("dbTags id->", dbTags[i].ID, ", tag id:", device.Tags[j].ID)
			if dbTags[i].ID == device.Tags[j].ID {
				findFlag = true
				break
			}
		}

		// 删除掉关联数据
		if !findFlag {
			logger.Info("need delete record->", dbTags[i].MidId)
			err := models.DB().Delete(&tms.DeviceAndTagMid{
				Model: gorm.Model{ID: dbTags[i].MidId},
			}).Error

			if err != nil {
				logger.Warn("Delete fail->", err.Error())
				return conf.DBError
			}
		}
	}

	logger.Info("find need create record-> db len:", len(dbTags), ", tags len:", len(device.Tags))
	for i := 0; i < len(device.Tags); i++ {
		findFlag := false
		for j := 0; j < len(dbTags); j++ {
			if dbTags[i].ID == device.Tags[j].ID {
				findFlag = true
				break
			}
		}

		// 添加关联数据
		if !findFlag {
			logger.Info("need create record->", device.Tags[i].ID)
			err := models.DB().Create(&tms.DeviceAndTagMid{
				TagID:    device.Tags[i].ID,
				DeviceId: device.ID,
			}).Error

			if err != nil {
				logger.Warn("Create fail->", err.Error())
				return conf.DBError
			}
		}
	}

	return conf.SUCCESS
}
