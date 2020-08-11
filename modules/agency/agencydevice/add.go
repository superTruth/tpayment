package agencydevice

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/download"
	"tpayment/pkg/fileutils"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(DeviceBindRequest)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	if req.AgencyId == 0 || (req.DeviceId == 0 && req.FileUrl == "") {
		logger.Warn("ParameterError")
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	var handleRet conf.ResultCode
	if req.DeviceId != 0 {
		handleRet = AddByID(ctx, req.AgencyId, req.DeviceId)
	} else {
		handleRet = AddByFile(ctx, req.AgencyId, req.FileUrl)
	}

	if handleRet != conf.SUCCESS {
		modules.BaseError(ctx, handleRet)
		return errors.New(handleRet.String())
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}

// 单个添加设备
func AddByID(ctx echo.Context, agencyId, deviceId uint) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
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

	return conf.SUCCESS
}

// 批量文件添加设备
const downloadDir = "./agencydevicefiles/"

func AddByFile(ctx echo.Context, agencyId uint, fileUrl string) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

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
	f, err := os.Open(localFilePath)
	// nolint
	defer f.Close()
	if err != nil {
		logger.Warn("open file err->", err.Error())
		return conf.UnknownError
	}
	buf := bufio.NewReader(f)
	r := csv.NewReader(buf)

	for i := 0; ; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Warn("read file err->", err.Error())
			return conf.UnknownError
		}

		// 跳过空值
		if len(record) == 0 || len(record[0]) < 5 {
			continue
		}

		// 查询一下是否已经存在这个device id
		device, err := tms.GetDeviceBySn(models.DB(), ctx, record[0])
		if err != nil {
			logger.Error("GetDeviceBySn fail->", err.Error())
			return conf.DBError
		}
		// 已经存在的情况
		if device != nil {
			logger.Info("device exist->", device.DeviceSn)
			// 先判断agency id是否相同，如果相同直接跳过
			if device.AgencyId == agencyId {
				continue
			}

			err = models.UpdateBaseRecord(&tms.DeviceInfo{
				BaseModel: models.BaseModel{
					ID: device.ID,
				},
				AgencyId: agencyId,
			})

			if err != nil {
				logger.Error("Updates fail->", err.Error())
				return conf.DBError
			}
			continue
		}
		logger.Info("device not exist->", record[0])
		// 如果不存在的情况，需要新建数据
		newDevice := tms.GenerateDeviceInfo()
		newDevice.AgencyId = agencyId
		newDevice.DeviceSn = record[0]
		err = models.CreateBaseRecord(newDevice)

		if err != nil {
			logger.Error("Create fail->", err.Error())
			return conf.DBError
		}
	}

	return conf.SUCCESS
}
