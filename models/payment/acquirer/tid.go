package acquirer

import (
	"math/rand"
	"time"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/pkg/tlog"

	"github.com/jinzhu/gorm"
)

type Terminal struct {
	models.BaseModel
	DeviceID          string `gorm:"column:device_id"`
	MerchantAccountID uint   `gorm:"column:merchant_account_id"`
	TID               string `gorm:"column:tid"`
	Addition          string `gorm:"column:addition"`
	AvailableAt       int64  `gorm:"column:available_at"`
}

func (Terminal) TableName() string {
	return "payment_tid"
}

// 获取
func (tid *Terminal) Get(id uint) (*Terminal, error) {
	ret := new(Terminal)
	err := tid.Db.Model(tid).Where("id=?", id).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

// 获取一条可用的TID
func (tid *Terminal) GetOneAvailable(merchantAccountID uint, deviceID string) (*Terminal, conf.ResultCode) {
	var (
		ret      = new(Terminal)
		freeTIDs []*Terminal
		err      error
	)
	logger := tlog.GetLogger(tid.Ctx)

	nowTimeStamp := time.Now().Unix()
	// 先查找是否有绑定这台设备的TID
	if deviceID != "" {
		err = tid.Db.Model(tid).Where("merchant_account_id=? AND device_id=?",
			merchantAccountID, deviceID).First(ret).Error
		if err != nil {
			if gorm.ErrRecordNotFound != err {
				logger.Error("find tid bind device fail->", err)
				return nil, conf.DBError
			}
			// 没有记录，则需要分配一条可用的TID
		} else { // 正常查到，确认是否可用
			if ret.AvailableAt > nowTimeStamp { // nolint
				logger.Warn("tid is busy:", "account id->", merchantAccountID, "deviceID->", deviceID)
				return nil, conf.TIDIsBusy
			}
			return ret, conf.Success
		}
	}

	// 分配一条公共TID给使用
	err = tid.Db.Model(tid).Where("merchant_account_id=? AND device_id=? AND available_at<?",
		merchantAccountID, deviceID, nowTimeStamp).Find(&freeTIDs).Error
	if err != nil {
		return nil, conf.DBError
	}

	if len(freeTIDs) == 0 {
		logger.Warn(conf.NoAvailableTID.String(), "->", "account id->", merchantAccountID, "deviceID->", deviceID)
		return nil, conf.NoAvailableTID
	}

	// 随机分配一条记录出去
	return freeTIDs[rand.Int()%len(freeTIDs)], conf.Success
}

// 查看一共有多少TID
func (tid *Terminal) GetTotal(merchantAccountID uint) (int, error) {
	ret := 0
	err := tid.Db.Model(tid).Where("merchant_account_id=?",
		merchantAccountID).Count(&ret).Error
	if err != nil {
		return 0, err
	}

	return ret, nil
}

// 锁定TID
func (tid *Terminal) Lock(timeOut time.Duration) conf.ResultCode {
	logger := tlog.GetLogger(tid.Ctx)

	expTime := time.Now().Unix() + int64(timeOut/time.Second)
	db := tid.Db.Model(tid).
		Where("id=? AND available_at=?", tid.ID, tid.AvailableAt).
		Update("available_at", expTime)

	err := db.Error
	if err != nil {
		logger.Error("lock tid fail(update)->", err.Error())
		return conf.DBError
	}

	if db.RowsAffected == 0 {
		logger.Warn("tid ", tid.ID, "is busy")
		return conf.TIDIsBusy
	}
	tid.AvailableAt = expTime
	return conf.Success
}

// 解锁TID
func (tid *Terminal) UnLock() conf.ResultCode {
	logger := tlog.GetLogger(tid.Ctx)

	db := tid.Db.Model(tid).
		Where("id=? AND available_at=?", tid.ID, tid.AvailableAt).
		Update("available_at", 0)

	err := db.Error
	if err != nil {
		logger.Error("unlock tid fail(update)->", err.Error())
		return conf.DBError
	}

	if db.RowsAffected == 0 {
		logger.Warn("tid ", tid.ID, "is busy")
		return conf.TIDIsBusy
	}
	tid.AvailableAt = 0
	return conf.Success
}
