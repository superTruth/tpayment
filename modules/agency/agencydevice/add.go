package agencydevice

import (
	"strings"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/download"
	"tpayment/pkg/fileutils"
	"tpayment/pkg/goroutine"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/xuri/excelize/v2"

	"github.com/gin-gonic/gin"
)

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetGoroutineLogger()

	req := new(DeviceBindRequest)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	if req.AgencyId == 0 || (req.DeviceId == 0 && req.FileUrl == "") {
		logger.Warn("ParameterError")
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	var handleRet conf.ResultCode
	if req.DeviceId != 0 {
		handleRet = AddByID(req.AgencyId, req.DeviceId)
	} else {
		handleRet = conf.Success
		goroutine.Go(func() {
			tlog.SetGoroutineLogger(logger) // 切换协程，承接log
			_ = AddByFile(req.AgencyId, req.FileUrl)
		})
	}

	if handleRet != conf.Success {
		modules.BaseError(ctx, handleRet)
		return
	}

	modules.BaseSuccess(ctx, nil)
}

// 单个添加设备
func AddByID(agencyId, deviceId uint64) conf.ResultCode {
	logger := tlog.GetGoroutineLogger()
	device := tms.DeviceInfo{
		BaseModel: models.BaseModel{
			ID: deviceId,
		},
		AgencyId: agencyId,
	}

	err := models.UpdateBaseRecord(device)
	if err != nil {
		logger.Info("UpdateBaseRecord sql error->", err.Error())
		return conf.DBError
	}

	return conf.Success
}

// 批量文件添加设备
const downloadDir = "./agencydevicefiles/"

func AddByFile(agencyId uint64, fileUrl string) conf.ResultCode {
	logger := tlog.GetGoroutineLogger()

	// 先下载文件
	_, fileName, _ := fileutils.SeparateFilePath(fileUrl)
	localFilePath := downloadDir + fileName
	err := download.Download(fileUrl, localFilePath)
	if err != nil {
		logger.Warn("download fail->", err.Error())
		return conf.UnknownError
	}

	// nolint
	defer fileutils.DeleteFile(localFilePath)

	// 读取里面的数据
	f, err := excelize.OpenFile(localFilePath)
	if err != nil {
		logger.Warn("read file fail->", err.Error())
		return conf.UnknownError
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		logger.Warn("can not find data from sheet1:", err.Error())
		return conf.UnknownError
	}

	// 提前读取出所有tag，防止后面一直读取
	tagArray, err := tms.DeviceTagDao.GetInAgency(agencyId)
	if err != nil {
		logger.Warn("get in agency err->", err.Error())
		return conf.DBError
	}
	tagMap := make(map[string]*tms.DeviceTag, len(tagArray)) // 转换成map，后面好做匹配
	for i := 0; i < len(tagArray); i++ {
		tagMap[tagArray[i].Name] = tagArray[i]
	}

	for i := 1; i < len(rows); i++ {
		record := rows[i]

		// 跳过空值
		if len(record) == 0 || len(record[0]) < 5 {
			continue
		}

		// 处理设备类型
		modelID, _ := handleDeviceModel(record[2])

		// 处理device
		device, handleRet := handleDevice(record[0], agencyId, modelID)
		if handleRet != conf.Success {
			return handleRet
		}

		// 处理tag
		if len(record) < 2 {
			logger.Info("device ", device.DeviceSn, " no tag, skip it")
			continue
		}

		handleRet = handleTag(record[1], device, agencyId, &tagMap)
		if handleRet != conf.Success {
			return handleRet
		}
	}

	return conf.Success
}

func handleDeviceModel(deviceModel string) (uint64, error) {
	//
	model, err := tms.GetModelByName(deviceModel)
	if err != nil {
		return 0, err
	}
	if model != nil {
		return model.ID, nil
	}

	// 如果不存在，就需要创建
	newModel := &tms.DeviceModel{
		Name: deviceModel,
	}
	err = models.CreateBaseRecord(newModel)
	if err != nil {
		return 0, err
	}

	return newModel.ID, nil
}

func handleDevice(deviceSn string, agencyId, modelID uint64) (*tms.DeviceInfo, conf.ResultCode) {
	logger := tlog.GetGoroutineLogger()

	// 查询一下是否已经存在这个device id
	device, err := tms.GetDeviceBySn(deviceSn)
	if err != nil {
		logger.Error("GetDeviceBySn fail->", err.Error())
		return nil, conf.DBError
	}
	// 已经存在的情况
	if device != nil {
		logger.Info("device exist->", device.DeviceSn)
		// 先判断agency id是否相同，如果相同直接跳过
		if device.AgencyId == agencyId {
			return device, conf.Success
		}

		err = models.UpdateBaseRecord(&tms.DeviceInfo{
			BaseModel: models.BaseModel{
				ID: device.ID,
			},
			AgencyId: agencyId,
		})

		if err != nil {
			logger.Error("Updates fail->", err.Error())
			return nil, conf.DBError
		}
		return device, conf.Success
	}
	logger.Info("device not exist->", deviceSn)
	// 如果不存在的情况，需要新建数据
	newDevice := tms.GenerateDeviceInfo()
	newDevice.AgencyId = agencyId
	newDevice.DeviceSn = deviceSn
	newDevice.DeviceModel = modelID
	err = models.CreateBaseRecord(newDevice)

	if err != nil {
		logger.Error("Create fail->", err.Error())
		return nil, conf.DBError
	}

	return newDevice, conf.Success
}

func handleTag(tagsDest string, device *tms.DeviceInfo, agencyId uint64, tags *map[string]*tms.DeviceTag) conf.ResultCode {
	logger := tlog.GetGoroutineLogger()

	// 删除掉现有所有tag
	err := tms.DeviceAndTagMidDao.DeleteAllTags(device.ID)
	if err != nil {
		logger.Error("DeviceTagDao.DeleteAllTags error->", err.Error())
		return conf.DBError
	}

	// 创建新tag
	tagsDestArray := strings.Split(tagsDest, ",")
	if len(tagsDestArray) == 0 {
		logger.Infof("no tag need add to the device")
		return conf.Success
	}

	for _, tag := range tagsDestArray {
		// 空字符串去除掉
		tmpTag := strings.ReplaceAll(tag, " ", "")
		if tmpTag == "" {
			continue
		}

		destTag, ok := (*tags)[tag]
		if !ok { // 如果不存在的tag，则添加一下
			destTag = &tms.DeviceTag{
				AgencyId: agencyId,
				Name:     tag,
			}

			err := tms.DeviceTagDao.Create(destTag)
			if err != nil {
				logger.Error("DeviceTagDao.Create error->", err.Error())
				return conf.DBError
			}
			(*tags)[destTag.Name] = destTag
		}

		err = tms.DeviceAndTagMidDao.Create(&tms.DeviceAndTagMid{
			TagID:    destTag.ID,
			DeviceId: device.ID,
		})

		if err != nil {
			logger.Error("create tag mid err ->", err.Error())
			return conf.DBError
		}
	}

	return conf.Success
}
