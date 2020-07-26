package uploadfile

import (
	"github.com/labstack/echo"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"
)

// AWSAccessKeyId=AKIAJPG36DFOEJ3XLTCA
// AWSSecretKey=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
// TODO 未完成
func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(tms.FileUpload)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// TODO 未做判断：当前用户可能没有此机构权限
	bean, err := tms.GetUploadFileByID(models.DB(), ctx, req.ID)
	if err != nil {
		logger.Error("GetAppByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	if bean == nil {
		logger.Info("GetAppByID sql error->", err.Error())
		modules.BaseError(ctx, conf.RecordNotFund)
		return err
	}

	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}
