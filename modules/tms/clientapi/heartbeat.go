package clientapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/goroutine"
	"tpayment/pkg/tlog"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

/**
数据同步
*/
var once sync.Once

func HearBeat(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	once.Do(func() {
		readDeviceModels(ctx)
	}) // 启动一次读取设备类型，并且后面每10分钟同步一次，减少没必要的IO操作

	data, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Warn("RequestApprove ReadAll->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	logger.Info("body->", string(data))

	bean := new(RequestBean)
	if err := json.Unmarshal(data, bean); err != nil {
		logger.Warn("Unmarshal", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 检查参数是否正确
	if bean.DeviceSn == "" || bean.DeviceModel == "" {
		logger.Warn("parameters miss")
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 查询出当前设备的信息
	deviceInfo, err := tms.GetDeviceBySn(models.DB(), ctx, bean.DeviceSn)
	if err != nil {
		logger.Info("GetDeviceBySn sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}
	if deviceInfo == nil { // 设备未创建的情况，需要创建
		logger.Info("设备未创建的情况，需要创建")
		deviceInfo = tms.GenerateDeviceInfo()
		copyRequestInfo2DeviceInfo(bean, deviceInfo) // 新上送的数据覆盖旧数据
		err = models.DB().Create(deviceInfo).Error
	} else {
		logger.Info("设备已经创建的情况")
		copyRequestInfo2DeviceInfo(bean, deviceInfo) // 新上送的数据覆盖旧数据
		err = models.DB().Update(deviceInfo).Error
	}
	if err != nil {
		logger.Error("设备更新或者创建失败->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	// 查询当前设备的app信息
	const PageLen = 1000

	requestApps := bean.AppInfos
	var retApps []AppInfo
	for i := 0; ; i++ {
		logger.Info("查询一次记录->", i)
		dbApps, err := tms.GetAppsInDevice(models.DB(), ctx, deviceInfo.ID, tms.AppInDeviceExternalIdTypeDevice, i*PageLen, PageLen) // 最多查出200条记录
		if err != nil {
			logger.Error("GetAppsInDevice error->", err.Error())
			modules.BaseError(ctx, conf.DBError)
			return err
		}
		logger.Info("最多查出200条记录->", len(dbApps))

		// 每次compare后，会返回requestApps变量，用来表示未在当前查询批次匹配到的记录，用于下次查询批次再次匹配
		var returnTmp []AppInfo
		logger.Info("开始对比数据==============")
		var errorCode conf.ResultCode
		requestApps, returnTmp, errorCode = compareApps(ctx, requestApps, dbApps)
		logger.Info("结束对比数据==============")
		if err != nil {
			modules.BaseError(ctx, errorCode)
			return err
		}
		for i := 0; i < len(returnTmp); i++ {
			retApps = append(retApps, returnTmp[i])
		}

		// 数据已经遍历完毕
		if len(dbApps) < PageLen {
			logger.Info("数据库遍历完毕")
			break
		}
	}

	// 未遍历到的数据，需要添加进数据库
	for i := 0; i < len(requestApps); i++ {
		appInDevice := new(tms.AppInDevice)
		appInDevice.Status = conf.TmsStatusWarningInstalled
		appInDevice.Name = requestApps[i].Name
		appInDevice.PackageId = requestApps[i].PackageId
		appInDevice.ExternalId = deviceInfo.ID
		appInDevice.VersionName = requestApps[i].VersionName
		appInDevice.VersionCode = requestApps[i].VersionCode
		appInDevice.ExternalIdType = tms.AppInDeviceExternalIdTypeDevice

		models.DB().Create(appInDevice)
	}

	//BaseSuccess(context)

	ret := new(ResponseBean)
	if deviceInfo.DeviceCsn != "" {
		ret.DeviceCsn = deviceInfo.DeviceCsn
	}
	if deviceInfo.RebootMode != "" {
		ret.RebootMode = deviceInfo.RebootMode
	}
	if deviceInfo.Alias != "" {
		ret.Alias = deviceInfo.Alias
	}
	if deviceInfo.RebootTime != "" {
		ret.RebootTime = deviceInfo.RebootTime
	}
	if deviceInfo.RebootDayInMonth != 0 {
		ret.RebootDayInMonth = deviceInfo.RebootDayInMonth
	}
	if deviceInfo.RebootDayInWeek != 0 {
		ret.RebootDayInWeek = deviceInfo.RebootDayInWeek
	}
	ret.AppInfos = retApps
	_ = ctx.JSON(http.StatusOK, ret)

	return nil
}

// 匹配处理历史记录，并且返回需要前端处理的数据
/*
	输入参数：requestApps： 需要匹配的数据
			dbApps：数据库查询出来的数据
	输出参数：
		1. 请求数据里面，未在数据库里面匹配成功的数据，需要下次继续匹配
		2. 匹配成功，需要返回给客户端的数据
		3. 错误信息
*/
// 返回数据：
func compareApps(ctx echo.Context, requestApps []AppInfo, dbApps []tms.AppInDevice) ([]AppInfo, []AppInfo, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)

	// 未匹配到的上送数据
	var unmatchApps []AppInfo

	var needReturnApp []AppInfo

	// 数组转map
	dbAppsMap := make(map[string]*tms.AppInDevice)
	for i := 0; i < len(dbApps); i++ {
		if dbApps[i].App == nil { // 如果数据库里面没有配置，则直接使用缓存的package id
			if dbApps[i].PackageId == "" {
				continue
			}
			dbAppsMap[dbApps[i].PackageId] = &dbApps[i]
		} else {
			if dbApps[i].App.PackageId == "" {
				continue
			}
			dbAppsMap[dbApps[i].App.PackageId] = &dbApps[i]
		}
	}

	for i := 0; i < len(requestApps); i++ {
		dbApp := dbAppsMap[requestApps[i].PackageId]
		if dbApp == nil { // 将匹配不到的数据返回，用于下次匹配
			logger.Info("将匹配不到的数据返回，用于下次匹配->", requestApps[i].PackageId)
			unmatchApps = append(unmatchApps, requestApps[i])
			continue
		} else { // 上报数据匹配到了配置数据
			logger.Info("上报数据匹配到了配置数据->", requestApps[i].PackageId)
			dbAppsMap[requestApps[i].PackageId] = nil // 去除已经匹配到的数据

			// 提前存储一下当前app状态, 做个判断，减少没必要的数据库操作
			if dbApp.VersionName == "" || dbApp.VersionCode == 0 || dbApp.PackageId == "" || dbApp.Name == "" ||
				dbApp.VersionName != requestApps[i].VersionName ||
				dbApp.VersionCode != requestApps[i].VersionCode ||
				dbApp.PackageId != requestApps[i].PackageId ||
				dbApp.Name != requestApps[i].Name {
				logger.Info("更新设备里面app状态->", requestApps[i].PackageId)
				dbApp.VersionName = requestApps[i].VersionName
				dbApp.VersionCode = requestApps[i].VersionCode
				dbApp.PackageId = requestApps[i].PackageId
				dbApp.Name = requestApps[i].Name
				models.DB().Updates(&dbApp)
			}

			if dbApp.Status == "" {
				continue
			}

			switch dbApp.Status {
			case conf.TmsStatusPendingInstall:
				if dbApp.AppFile == nil || dbApp.AppFile.VersionCode == 0 {
					logger.Warn("STATUS_PENDING_INSTALL VersionCode is nil->", requestApps[i].PackageId)
					continue
				}
				logger.Info("conf.STATUS_PENDING_INSTALL: current->", requestApps[i].VersionCode, ", config->", dbApp.AppFile.VersionCode)
				if requestApps[i].VersionCode < dbApp.AppFile.VersionCode { // 上报app version code小于当前，下发当前配置数据
					newApp := generateAppFromConfig(dbApp)

					newApp.Status = conf.TmsStatusPendingInstall
					needReturnApp = append(needReturnApp, *newApp)
					continue
				}

				if requestApps[i].VersionCode == dbApp.AppFile.VersionCode { // 上报app version code等于当前，变成Installed状态，并且不下发
					dbApp.Status = conf.TmsStatusInstalled
					models.DB().Updates(dbApp)
					continue
				}
				// 上报app version code大于当前配置，更新数据为上报的数据，并且不下发，状态改为Warning Installed
				dbApp.Status = conf.TmsStatusWarningInstalled
				models.DB().Updates(dbApp)
				continue

			case conf.TmsStatusInstalled:
				if dbApp.AppFile == nil || dbApp.AppFile.VersionCode == 0 {
					logger.Warn("STATUS_INSTALLED VersionCode is nil->", requestApps[i].PackageId)
					continue
				}
				logger.Info("conf.STATUS_INSTALLED", requestApps[i].VersionCode, ", config->", dbApp.AppFile.VersionCode)
				if requestApps[i].VersionCode < dbApp.AppFile.VersionCode { // 上报app version code小于当前，下发当前配置数据，状态改为Pending Install
					newApp := generateAppFromConfig(dbApp)
					dbApp.Status = conf.TmsStatusPendingInstall
					models.DB().Updates(dbApp)

					newApp.Status = conf.TmsStatusPendingInstall
					needReturnApp = append(needReturnApp, *newApp)
					continue
				}

				if requestApps[i].VersionCode == dbApp.AppFile.VersionCode { // 上报app version code等于当前，不下发
					continue
				}

				// 上报app version code大于当前配置，更新数据为上报的数据，并且不下发，状态改为Warning Installed
				dbApp.Status = conf.TmsStatusWarningInstalled
				models.DB().Updates(dbApp)

			case conf.TmsStatusPendingUninstalled:
				logger.Info("conf.STATUS_PENDING_UNINSTALL", requestApps[i].VersionCode)
				newApp := new(AppInfo)
				newApp.Name = requestApps[i].Name
				newApp.PackageId = requestApps[i].PackageId
				newApp.Status = conf.TmsStatusPendingUninstalled
				needReturnApp = append(needReturnApp, *newApp)

			case conf.TmsStatusWarningInstalled:
				logger.Info("conf.STATUS_WARNING_INSTALLED", requestApps[i].VersionCode)
				if dbApp.AppFile == nil || dbApp.AppFile.VersionCode == 0 { // 从来没有安装过的情况，不处理
					continue
				}

				if requestApps[i].VersionCode < dbApp.AppFile.VersionCode { // 上报app version code小于当前，下发当前配置数据，状态改为Pending Install
					newApp := generateAppFromConfig(dbApp)
					dbApp.Status = conf.TmsStatusPendingInstall
					models.DB().Updates(dbApp)

					newApp.Status = conf.TmsStatusPendingInstall
					needReturnApp = append(needReturnApp, *newApp)
					continue
				}

				if requestApps[i].VersionCode == dbApp.AppFile.VersionCode { // 等于当前，更新状态为Installed，不下发
					dbApp.Status = conf.TmsStatusInstalled
					models.DB().Updates(dbApp)
					continue
				}

				// 大于当前，不下发
			default:
				logger.Error("conf.default")
			}
		}
	}

	// 未上报的配置数据
	for _, dbApp := range dbAppsMap {
		if dbApp == nil || dbApp.Status == "" {
			continue
		}

		switch dbApp.Status {
		case conf.TmsStatusPendingInstall: // 下发当前配置数据
			newApp := generateAppFromConfig(dbApp)

			newApp.Status = conf.TmsStatusPendingInstall
			needReturnApp = append(needReturnApp, *newApp)
			continue
		case conf.TmsStatusInstalled: // 下发当前配置数据，并且状态改为Pending Install
			newApp := generateAppFromConfig(dbApp)

			newApp.Status = conf.TmsStatusPendingInstall
			needReturnApp = append(needReturnApp, *newApp)
			dbApp.Status = conf.TmsStatusPendingInstall
			models.DB().Updates(dbApp)
			continue
		case conf.TmsStatusPendingUninstalled: // 不下发配置，删除掉这个记录
			models.DB().Delete(&dbApp)

		case conf.TmsStatusWarningInstalled:
			if dbApp.AppFile == nil || dbApp.AppFile.VersionCode == 0 { // 从来没有配置过的情况，删除掉
				models.DB().Delete(&dbApp)
				continue
			}

			// 配置过的话，下发当前配置数据，并且状态改为Pending Install
			newApp := generateAppFromConfig(dbApp)

			newApp.Status = conf.TmsStatusPendingInstall
			needReturnApp = append(needReturnApp, *newApp)
			dbApp.Status = conf.TmsStatusPendingInstall
			models.DB().Updates(&dbApp)
			continue

		default:
		}
	}

	return unmatchApps, needReturnApp, conf.SUCCESS
}

// 拷贝请求设备信息到设备信息
func copyRequestInfo2DeviceInfo(requestDevice *RequestBean, deviceInfo *tms.DeviceInfo) {
	deviceInfo.DeviceSn = (requestDevice.DeviceSn)

	if requestDevice.DeviceCsn != "" {
		deviceInfo.DeviceCsn = (requestDevice.DeviceCsn)
	}

	if requestDevice.DeviceModel != "" && deviceModels != nil {
		// TODO
		deviceInfo.DeviceModel = requestDevice.DeviceModel
	}

	if requestDevice.LocationLat != "" {
		deviceInfo.LocationLat = (requestDevice.LocationLat)
	}

	if requestDevice.LocationLon != "" {
		deviceInfo.LocationLon = (requestDevice.LocationLon)
	}

	if requestDevice.PushToken != "" {
		deviceInfo.PushToken = (requestDevice.PushToken)
	}

	if requestDevice.Power != 0 {
		deviceInfo.Battery = requestDevice.Power
	}

}

// 从配置App数据生成升级信息
func generateAppFromConfig(configApp *tms.AppInDevice) *AppInfo {
	retApp := new(AppInfo)

	if configApp.App.Name != "" {
		retApp.Name = configApp.App.Name
	}
	if configApp.App.PackageId != "" {
		retApp.PackageId = configApp.App.PackageId
	}

	if configApp.AppFile.VersionName != "" {
		retApp.VersionName = configApp.AppFile.VersionName
	}
	if configApp.AppFile.VersionCode != 0 {
		retApp.VersionCode = configApp.AppFile.VersionCode
	}
	if configApp.AppFile.UpdateDescription != "" {
		retApp.Description = configApp.AppFile.UpdateDescription
	}

	retApp.FileInfo = new(FileInfo)
	if configApp.AppFile.FileName != "" {
		retApp.FileInfo.Name = configApp.AppFile.FileName
	}
	if configApp.AppFile.FileUrl != "" {
		retApp.FileInfo.Url = configApp.AppFile.FileUrl
	}

	return retApp
}

var deviceModels map[string]uint

func readDeviceModels(ctx echo.Context) {
	goroutine.Go(func() {
		logger := tlog.Logger{}
		logger.Init(uuid.New().String())
		for {
			modelArray, err := tms.GetModels()
			if err == nil {
				deviceModelTmp := make(map[string]uint)
				for _, v := range modelArray {
					deviceModelTmp[*v.Name] = v.ID
				}
				deviceModels = deviceModelTmp
			}
			logger.Info("读取一次models->", len(modelArray))
			time.Sleep(time.Minute * 10) // 10分钟同步一次
		}
	}, ctx)
}
