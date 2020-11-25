package device

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

func UpdateHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(tms.DeviceInfo)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 查询是否已经存在的账号
	bean, err := tms.GetDeviceByID(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Info("GetDeviceByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}
	if bean == nil {
		logger.Warn(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 判断权限
	if tms2.CheckPermission(ctx, bean) != nil {
		logger.Warn(conf.NoPermission.String())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	// 生成新账号
	err = models.UpdateBaseRecord(req)

	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	// 合并tags
	errorCode := mergeTags(ctx, req)
	if errorCode != conf.Success {
		modules.BaseError(ctx, errorCode)
		return
	}

	modules.BaseSuccess(ctx, nil)
}

// 合并tags
func mergeTags(ctx *gin.Context, device *tms.DeviceInfo) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	// 前端没传入
	if device.Tags == nil {
		return conf.Success
	}

	dbTags, err := tms.QueryTagsInDevice(models.DB(), ctx, device)

	if err != nil {
		logger.Warn("QueryTagsInDevice fail->", err.Error())
		return conf.DBError
	}

	// 先找到需要删除的tag(在数据库中有，但是update没有的)
	logger.Info("find need delete record-> db len:", len(dbTags), ", tags len:", len(*device.Tags))
	for i := 0; i < len(dbTags); i++ {
		findFlag := false

		for j := 0; j < len(*device.Tags); j++ {
			logger.Info("dbTags id->", dbTags[i].ID, ", tag id:", (*device.Tags)[j].ID)
			if dbTags[i].ID == (*device.Tags)[j].ID {
				findFlag = true
				break
			}
		}

		// 删除掉关联数据
		if !findFlag {
			logger.Info("need delete record->", dbTags[i].MidId)
			err := models.DB().Delete(&tms.DeviceAndTagMid{
				BaseModel: models.BaseModel{ID: dbTags[i].MidId},
			}).Error

			if err != nil {
				logger.Warn("Delete fail->", err.Error())
				return conf.DBError
			}
		}
	}

	logger.Info("find need create record-> db len:", len(dbTags), ", tags len:", len(*device.Tags))
	for i := 0; i < len(*device.Tags); i++ {
		findFlag := false
		for j := 0; j < len(dbTags); j++ {
			if dbTags[j].ID == (*device.Tags)[i].ID {
				findFlag = true
				break
			}
		}

		// 添加关联数据
		if !findFlag {
			logger.Info("need create record->", (*device.Tags)[i].ID)
			err := models.DB().Create(&tms.DeviceAndTagMid{
				TagID:    (*device.Tags)[i].ID,
				DeviceId: device.ID,
			}).Error

			if err != nil {
				logger.Warn("Create fail->", err.Error())
				return conf.DBError
			}
		}
	}

	return conf.Success
}
