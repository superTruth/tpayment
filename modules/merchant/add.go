package merchant

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"tpayment/conf"
	"tpayment/internal/basekey"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/agency"
	"tpayment/models/merchant"
	"tpayment/models/tms"
	"tpayment/modules"
	"tpayment/pkg/download"
	"tpayment/pkg/fileutils"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AddHandle(ctx *gin.Context) {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.Merchant)

	err := utils.Body2Json(ctx.Request.Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	// 判断当前agency id
	req.AgencyId, err = modules.GetAgencyId(ctx, req.AgencyId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return
	}

	var retCode conf.ResultCode
	if req.FileUrl == "" {
		retCode = addNormal(ctx, req)
	} else {
		retCode = addByFile(ctx, req)
	}

	if retCode != conf.Success {
		modules.BaseError(ctx, retCode)
		return
	}
	modules.BaseSuccess(ctx, nil)
}

// 常规添加
func addNormal(ctx *gin.Context, req *merchant.Merchant) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var err error

	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Error("CreateBaseRecord sql error->", err.Error())
		return conf.DBError
	}

	return conf.Success
}

// 文件导入的方式
const downloadDir = "./merchant_file/"

type fileItemBean struct {
	MerchantName   string
	MerchantAddr   string
	MerchantTel    string
	MerchantEmail  string
	DeviceSN       string
	PaymentMethods *models.StringArray
	EntryTypes     *models.StringArray
	PaymentTypes   *models.StringArray
	BankName       string
	MID            string
	TID            string
	Addition       string
	StuffName      string
	StuffEmail     string
	StuffPwd       string
	StuffRole      string
}

func addByFile(ctx *gin.Context, req *merchant.Merchant) conf.ResultCode {
	logger := tlog.GetLogger(ctx)

	// 先下载文件
	_, fileName, _ := fileutils.SeparateFilePath(req.FileUrl)
	localFilePath := downloadDir + fileName
	err := download.Download(req.FileUrl, localFilePath)
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
	_, err = r.Read() // 跳过抬头
	if err != nil {
		logger.Warn("skip first row error->", err.Error())
		return conf.UnknownError
	}

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
		if len(record) < 15 {
			continue
		}

		fileItem := &fileItemBean{
			MerchantName:   record[0],
			MerchantAddr:   record[1],
			MerchantTel:    record[2],
			MerchantEmail:  record[3],
			DeviceSN:       record[4],
			PaymentMethods: convertArray(record[5]),
			EntryTypes:     convertArray(record[6]),
			PaymentTypes:   convertArray(record[7]),
			BankName:       record[8],
			MID:            record[9],
			TID:            record[10],
			Addition:       record[11],
			StuffName:      record[12],
			StuffEmail:     record[13],
			StuffPwd:       record[14],
			StuffRole:      record[15],
		}

		_ = handleFileItem(ctx, req.AgencyId, fileItem)
	}

	return conf.Success
}

func convertArray(src string) *models.StringArray {
	srcFormat := strings.ReplaceAll(src, " ", "")
	splitArray := strings.Split(srcFormat, ",")
	if len(splitArray) == 0 {
		return nil
	}
	ret := make(models.StringArray, 0)
	ret = append(ret, splitArray...)
	return &ret
}

func handleFileItem(ctx *gin.Context, agencyID uint64, fileItem *fileItemBean) error {
	log := tlog.GetLogger(ctx)

	// 添加商户
	merchantBean, err := merchant.Dao.GetByName(fileItem.MerchantName)
	if err != nil {
		log.Errorf("GetByName fail: %s", err.Error())
		return err
	}
	if merchantBean == nil { // 不存在的话，就新建一个
		merchantBean = &merchant.Merchant{
			AgencyId: agencyID,
			Name:     fileItem.MerchantName,
			Tel:      fileItem.MerchantTel,
			Addr:     fileItem.MerchantAddr,
			Email:    fileItem.MerchantEmail,
		}
		if err = models.CreateBaseRecord(merchantBean); err != nil {
			log.Errorf("create merchant fail: %s", err.Error())
			return err
		}
	}

	// 添加设备关联
	// 查找设备是否存在
	deviceInfo, err := tms.DeviceInfoDao.GetBySn(fileItem.DeviceSN)
	if err != nil {
		log.Errorf("get device sn fail: %s", err.Error())
		return err
	}
	if deviceInfo == nil {
		return fmt.Errorf("can not find the device: %s", fileItem.DeviceSN)
	}
	// 查看设备关联是否已经存在
	deviceInMerchant, err := merchant.DeviceInMerchantDao.GetByMerchantIdAndDeviceID(merchantBean.ID, deviceInfo.ID)
	if err != nil {
		log.Errorf("GetByMerchantIdAndDeviceID fail: %s", err.Error())
		return err
	}
	if deviceInMerchant == nil { // 创建
		deviceInMerchant = &merchant.DeviceInMerchant{
			DeviceId:   deviceInfo.ID,
			MerchantId: merchantBean.ID,
		}
		if err = models.CreateBaseRecord(deviceInMerchant); err != nil {
			log.Errorf("create device in merchant fail: %s", err.Error())
			return err
		}
	}

	// 查看银行ID
	acq, err := agency.AcquirerDao.GetByName(agencyID, fileItem.BankName)
	if err != nil {
		log.Errorf("get acquirer fail: %s", err.Error())
		return err
	}
	if acq == nil {
		return fmt.Errorf("can not find acq :%s", fileItem.BankName)
	}

	// 添加支付参数
	paymentSetting, err := merchant.PaymentSettingDao.GetByMidTid(deviceInMerchant.ID, fileItem.MID, fileItem.TID)
	if err != nil {
		log.Errorf("GetByMidTid fail: %s", err.Error())
		return err
	}
	if paymentSetting == nil {
		paymentSetting = &merchant.PaymentSettingInDevice{
			MerchantDeviceId: deviceInMerchant.ID,
			PaymentMethods:   fileItem.PaymentMethods,
			EntryTypes:       fileItem.EntryTypes,
			PaymentTypes:     fileItem.PaymentTypes,
			AcquirerId:       acq.ID,
			Mid:              fileItem.MID,
			Tid:              fileItem.TID,
			Addition:         fileItem.Addition,
		}

		err = models.CreateBaseRecord(paymentSetting)
		if err != nil {
			return fmt.Errorf("create acq fail: %s", err.Error())
		}
	}

	// 添加用户
	user, err := account.UserBeanDao.GetByEmail(agencyID, fileItem.StuffEmail)
	if err != nil {
		return fmt.Errorf("get user by email fail: %s", err.Error())
	}
	if user == nil {
		user = &account.UserBean{
			AgencyId: agencyID,
			Email:    fileItem.StuffEmail,
			Pwd:      basekey.Hash([]byte(fileItem.StuffPwd)),
			Name:     fileItem.StuffName,
			Role:     string(conf.RoleUser),
			Active:   true,
		}

		if err = models.CreateBaseRecord(user); err != nil {
			return fmt.Errorf("create user fail: %s", err.Error())
		}
	}
	// 查看user是否已经关联上去
	userMerchantAssociate, err := merchant.UserMerchantAssociateDao.GetByMerchantIdAndUserId(merchantBean.ID, user.ID)
	if err != nil {
		return fmt.Errorf("GetByMerchantIdAndUserId fail: %s", err.Error())
	}
	if userMerchantAssociate == nil { // 创建
		userMerchantAssociate = &merchant.UserMerchantAssociate{
			MerchantId: merchantBean.ID,
			UserId:     user.ID,
			Role:       fileItem.StuffRole,
		}
		if err = models.CreateBaseRecord(userMerchantAssociate); err != nil {
			return fmt.Errorf("create user associate fail: %s", err.Error())
		}
	}

	return nil
}
