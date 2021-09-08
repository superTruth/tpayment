package acquirer

import (
	"math/rand"
	"time"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/pkg/tlog"

	"gorm.io/gorm"
)

var TerminalDao = &Terminal{
	BaseModel: models.BaseModel{},
}

type Terminal struct {
	models.BaseModel
	DeviceID          string `gorm:"column:device_id"`
	MerchantAccountID uint64 `gorm:"column:merchant_account_id"`
	TID               string `gorm:"column:tid"`
	Addition          string `gorm:"column:addition"`
	AvailableAt       int64  `gorm:"column:available_at"`

	TraceNum uint64 `gorm:"column:trace_num"`
	BatchNum uint64 `gorm:"column:batch_num"`
}

func (Terminal) TableName() string {
	return "payment_tid"
}

// 获取
func (t *Terminal) Get(id uint64) (*Terminal, error) {
	ret := new(Terminal)
	err := models.DB.Model(t).Where("id=?", id).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

//
func (t *Terminal) GetByTID(mid uint64, tid string) (*Terminal, error) {
	ret := new(Terminal)
	err := models.DB.Model(t).Where("merchant_account_id=? and device_id=?", mid, tid).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return nil, nil
		}
		return nil, err
	}

	return ret, nil
}

// 获取一条可用的TID
func (t *Terminal) GetOneAvailable(merchantAccountID uint64, deviceID string) (*Terminal, conf.ResultCode) {
	var (
		ret      = new(Terminal)
		freeTIDs []*Terminal
		err      error
	)
	logger := tlog.GetGoroutineLogger()

	nowTimeStamp := time.Now().Unix()
	// 先查找是否有绑定这台设备的TID
	if deviceID != "" {
		err = models.DB.Model(t).Where("merchant_account_id=? AND device_id=?",
			merchantAccountID, deviceID).First(ret).Error
		if err != nil {
			if gorm.ErrRecordNotFound != err {
				logger.Error("find t bind device fail->", err)
				return nil, conf.DBError
			}
			// 没有记录，则需要分配一条可用的TID
		} else { // 正常查到，确认是否可用
			if ret.AvailableAt > nowTimeStamp { // nolint
				logger.Warn("t is busy:", "account id->", merchantAccountID, "deviceID->", deviceID)
				return nil, conf.TIDIsBusy
			}
			return ret, conf.Success
		}
	}

	// 分配一条公共TID给使用，公共TID是指device id数据为""
	err = models.DB.Model(t).Where("merchant_account_id=? AND available_at<? AND (device_id='' OR device_id IS NULL)",
		merchantAccountID, nowTimeStamp).Find(&freeTIDs).Error
	if err != nil {
		return nil, conf.DBError
	}

	if len(freeTIDs) == 0 {
		logger.Warn(conf.NoAvailableTID.String(), "->", "account id->", merchantAccountID, "deviceID->", deviceID)
		return nil, conf.NoAvailableTID
	}

	randIndex := rand.Int() % len(freeTIDs)
	logger.Info("random t ", freeTIDs[randIndex].TID, "from pool size->", len(freeTIDs))

	// 随机分配一条记录出去
	return freeTIDs[randIndex], conf.Success
}

// 获取MID下面的所有TID
func (t *Terminal) GetByMID(merchantAccountID uint64) ([]*Terminal, error) {
	var (
		tids []*Terminal
		err  error
	)

	err = models.DB.Model(t).Where("merchant_account_id=?",
		merchantAccountID).Find(&tids).Error
	if err != nil {
		return nil, err
	}

	return tids, nil
}

// 查看一共有多少TID
func (t *Terminal) GetTotal(merchantAccountID uint64) (int, error) {
	ret := int64(0)
	err := models.DB.Model(t).Where("merchant_account_id=?",
		merchantAccountID).Count(&ret).Error
	if err != nil {
		return 0, err
	}

	return int(ret), nil
}

// 锁定TID
func (t *Terminal) Lock(timeOut time.Duration) conf.ResultCode {
	logger := tlog.GetGoroutineLogger()
	timeNow := time.Now().Unix()
	if t.AvailableAt > timeNow {
		return conf.TIDIsBusy
	}

	expTime := timeNow + int64(timeOut/time.Second)
	db := models.DB.Model(t).
		Where("id=? AND available_at=?", t.ID, t.AvailableAt).
		Update("available_at", expTime)

	err := models.DB.Error
	if err != nil {
		logger.Error("lock t fail(update)->", err.Error())
		return conf.DBError
	}

	if db.RowsAffected == 0 {
		logger.Warn("t ", t.ID, "is busy")
		return conf.TIDIsBusy
	}
	t.AvailableAt = expTime
	return conf.Success
}

// 解锁TID
func (t *Terminal) UnLock() conf.ResultCode {
	logger := tlog.GetGoroutineLogger()

	db := models.DB.Model(t).
		Where("id=? AND available_at=?", t.ID, t.AvailableAt).
		Update("available_at", 0)

	err := models.DB.Error
	if err != nil {
		logger.Error("unlock t fail(update)->", err.Error())
		return conf.DBError
	}

	if db.RowsAffected == 0 {
		logger.Warn("t ", t.ID, "is busy")
		return conf.TIDIsBusy
	}
	t.AvailableAt = 0
	return conf.Success
}

func (t *Terminal) IncTraceNum() error {
	t.TraceNum = (t.TraceNum + 1) % 1000000
	if t.TraceNum == 0 {
		t.TraceNum = 1
	}

	return models.DB.Model(t).UpdateColumns(map[string]interface{}{"trace_num": t.TraceNum}).Error
}

func (t *Terminal) IncBatchNum() error {
	t.BatchNum = (t.BatchNum + 1) % 1000000
	if t.BatchNum == 0 {
		t.BatchNum = 1
	}

	return models.DB.Model(t).UpdateColumns(map[string]interface{}{"batch_num": t.BatchNum}).Error
}
