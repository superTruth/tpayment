package acquirer

import (
	"time"
	"tpayment/models"
)

type TIDPool struct {
	models.BaseModel

	MerchantAccountID uint   `gorm:"tid"`
	TID               string `gorm:"tid"`
	ExpiredAt         int64  `gorm:"expired_at"`
}

func (TIDPool) TableName() string {
	return "tid_pool"
}

// 查找
const (
	MaxLimitTID = 20              // 最大一次性搜索出20条数据
	ExpTime     = time.Minute * 5 // 5分钟的过期时间
)

func (pool *TIDPool) GetOneAvailable(acquirerID, MID string) (*TIDPool, error) {
	var ret []*TIDPool
	err := pool.Db.Model(pool).Where("acquirer_id=? AND mid=? AND expired_at<?",
		acquirerID, MID, time.Now().Unix()).Find(&ret).Error
	if err != nil {
		return nil, err
	}

	// 没有找到
	if len(ret) == 0 {
		return nil, nil
	}

	// 从找到的数据里面随机出来一条

	return ret[0], nil
}
