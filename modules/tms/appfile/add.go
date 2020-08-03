package appfile

import (
	"path"
	"regexp"
	"strings"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/apkparser"
	"tpayment/pkg/goroutine"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(tms.AppFile)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// TODO 未做判断：当前用户可能没有此机构权限
	bean, err := tms.GetAppFileByID(models.DB(), ctx, req.ID)
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

	req.Status = conf.AppFileStatusPending
	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	// 异步解析
	goroutine.Go(func() {
		StartDecode(ctx, req.ID)
	}, ctx)

	return nil
}

func StartDecode(ctx echo.Context, id uint) {
	logger := tlog.GetLogger(ctx)

	appFile, err := tms.GetAppFileByID(models.DB(), ctx, id)
	if err != nil {
		logger.Error("GetAppFileByID sql error->", err.Error())
		return
	}

	if appFile == nil { // 找不到记录
		logger.Error("找不到apk文件记录->", id)
		return
	}

	// 获取文件所属的app，用来判断解析后的apk文件package id和app设定package id是相同的防止错误
	app, err := tms.GetAppByID(models.DB(), ctx, *appFile.AppId)
	if err != nil {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "Can't get parent app"
		models.DB().Updates(appFile)
		logger.Warn("GetAppByID error->", appFile.DecodeFailMsg)
		return
	}

	// 开始
	appFile.Status = conf.AppFileStatusDecoding
	err = models.DB().Updates(appFile).Error
	if err != nil {
		logger.Error("Updates appFile sql err->", err.Error())
		return
	}

	// 下载文件
	if appFile.FileUrl == "" {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "fileutils url is empty"
		models.DB().Updates(appFile)
		logger.Warn("Updates appFile error->", appFile.DecodeFailMsg)
		return
	}

	// 如果文件不是.apk结尾，则认为不是apk文件，不解析
	ret, err := regexp.MatchString(`\.apk$`, strings.ToLower(appFile.FileUrl))
	if err != nil || !ret {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "fileutils url is not apk fileutils"
		models.DB().Updates(appFile)
		logger.Warn("Updates appFile error->", appFile.DecodeFailMsg)
		return
	}

	// 下载并且解析
	apkParser := apkparser.ApkParser{Url: appFile.FileUrl}
	apkInfo, err := apkParser.DownloadApkInfo()
	if err != nil {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "apk parse fail->" + err.Error()
		models.DB().Updates(appFile)
		logger.Warn(appFile.DecodeFailMsg)
		return
	}

	// 判断package id是否和parent的相同
	if app.PackageId != apkInfo.Package {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "Package is not same as parent app: parent->" + app.PackageId + ", current ->" + apkInfo.Package
		models.DB().Updates(appFile)
		logger.Warn(appFile.DecodeFailMsg)
		return
	}

	// 存储解析结果
	appFile.Status = conf.AppFileStatusDone
	appFile.VersionCode = apkInfo.VersionCode
	appFile.VersionName = apkInfo.VersionName
	appFile.FileName = path.Base(appFile.FileUrl)

	models.DB().Updates(appFile)
}
