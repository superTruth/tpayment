package acquirer

import (
	"tpayment/models"
)

type Account struct {
	models.BaseModel

	Tag string `gorm:"column:tag"`
	//TraceNum uint64   `gorm:"column:trace_num"`
	//BatchNum uint64   `gorm:"column:batch_num"`
}

func (Account) TableName() string {
	return "payment_acquirer_account"
}

//func (account *Account) Get(tag string) (*Account, error) {
//	var ret = new(Account)
//	err := account.Db.Model(account).Where("tag=?", tag).First(ret).Error
//	if err != nil {
//		if gorm.ErrRecordNotFound == err { // 没有记录, 就创建一条记录
//			return nil, nil
//		}
//		return nil, err
//	}
//
//	return ret, nil
//}
//
//func (account *Account) GetOrCreate(tag string) (*Account, error) {
//	var ret = new(Account)
//	err := account.Db.Model(account).Where("tag=?", tag).First(ret).Error
//	if err != nil {
//		if gorm.ErrRecordNotFound == err { // 没有记录, 就创建一条记录
//			ret = new(Account)
//			ret.Tag = tag
//			ret.BatchNum = 1
//			ret.TraceNum = 1
//			err = account.Db.Model(account).Create(ret).Error
//			if err != nil {
//				return nil, err
//			}
//			return ret, err
//		}
//		return nil, err
//	}
//
//	return ret, nil
//}
//
//func (account *Account) IncTraceNum() error {
//	account.TraceNum = (account.TraceNum + 1) % 1000000
//	if account.TraceNum == 0 {
//		account.TraceNum = 1
//	}
//
//	return account.Db.Model(account).Update(map[string]interface{}{"trace_num": account.TraceNum}).Error
//}
//
//func (account *Account) IncBatchNum() error {
//	account.BatchNum = (account.BatchNum + 1) % 1000000
//	if account.BatchNum == 0 {
//		account.BatchNum = 1
//	}
//
//	return account.Db.Model(account).Update(map[string]interface{}{"batch_num": account.BatchNum}).Error
//}
