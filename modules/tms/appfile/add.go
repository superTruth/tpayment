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

	"github.com/gin-gonic/gin"
)

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(tms.AppFile)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}
	if req.AppId == 0 {
		logger.Warn(conf.ParameterError.String())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 获取机构ID，系统管理员为0
	agencyId, err := modules.GetAgencyId2(ctx)
	if err != nil {
		logger.Warn("GetAgencyId2->", err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	bean, err := tms.GetAppByID(models.DB(), ctx, req.AppId)
	if err != nil {
		logger.Error("GetAppByID sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	if bean == nil {
		logger.Info(conf.RecordNotFund.String())
		modules.BaseError(ctx, conf.RecordNotFund)
		return
	}

	// 无权限删除
	if agencyId != 0 && agencyId != bean.AgencyId {
		logger.Warn("current agency id is:", bean.AgencyId, ", your id:", agencyId)
		modules.BaseError(ctx, conf.NoPermission)
		return
	}

	req.ID = 0
	req.Status = conf.AppFileStatusPending
	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return
	}

	modules.BaseSuccess(ctx, nil)

	// 异步解析
	goroutine.Go(func() {
		StartDecode(ctx, req.ID)
	}, ctx)
}

func StartDecode(ctx *gin.Context, id uint64) {
	logger := tlog.GetLogger(ctx)

	logger.Info("start decode app file->", id)

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
	app, err := tms.GetAppByID(models.DB(), ctx, appFile.AppId)
	if err != nil {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "Can't get parent app"
		_ = models.UpdateBaseRecord(appFile)
		logger.Warn("GetAppByID error->", appFile.DecodeFailMsg)
		return
	}

	// 开始
	appFile.Status = conf.AppFileStatusDecoding
	err = models.UpdateBaseRecord(appFile)
	if err != nil {
		logger.Error("Updates appFile sql err->", err.Error())
		return
	}

	// 下载文件
	if appFile.FileUrl == "" {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "fileutils url is empty"
		_ = models.UpdateBaseRecord(appFile)
		logger.Warn("Updates appFile error->", appFile.DecodeFailMsg)
		return
	}

	// 如果文件不是.apk结尾，则认为不是apk文件，不解析
	ret, err := regexp.MatchString(`\.apk$`, strings.ToLower(appFile.FileUrl))
	if err != nil || !ret {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "fileutils url is not apk fileutils"
		_ = models.UpdateBaseRecord(appFile)
		logger.Warn("Updates appFile error->", appFile.DecodeFailMsg)
		return
	}

	// 下载并且解析
	apkParser := apkparser.ApkParser{Url: appFile.FileUrl}
	apkInfo, err := apkParser.DownloadApkInfo()
	if err != nil {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "apk parse fail->" + err.Error()
		_ = models.UpdateBaseRecord(appFile)
		logger.Warn(appFile.DecodeFailMsg)
		return
	}

	// 判断package id是否和parent的相同
	if app.PackageId != apkInfo.Package {
		appFile.Status = conf.AppFileStatusFail
		appFile.DecodeFailMsg = "Package is not same as parent app: parent->" + app.PackageId + ", current ->" + apkInfo.Package
		_ = models.UpdateBaseRecord(appFile)
		logger.Warn(appFile.DecodeFailMsg)
		return
	}

	// 存储解析结果
	appFile.Status = conf.AppFileStatusDone
	appFile.VersionCode = apkInfo.VersionCode
	appFile.VersionName = apkInfo.VersionName
	appFile.FileName = path.Base(appFile.FileUrl)

	_ = models.UpdateBaseRecord(appFile)
}
