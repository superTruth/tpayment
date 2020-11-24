package acquirer

import (
	"tpayment/models"

	"github.com/jinzhu/gorm"
)

type Account struct {
	models.BaseModel

	AcquirerID string `gorm:"acquirer_id"`
	Tag        string `gorm:"tag"`
	TraceNum   uint   `gorm:"trace_num"`
	BatchNum   uint   `gorm:"batch_num"`
}

func (Account) TableName() string {
	return "acquirer_account"
}

func (account *Account) GetOrCreate(acquirerID, tag string) (*Account, error) {
	var ret *Account
	err := account.Db.Model(account).Where("acquirer_id=? AND tag=?",
		acquirerID, tag).First(ret).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录, 就创建一条记录
			ret = new(Account)
			ret.AcquirerID = acquirerID
			ret.Tag = tag
			ret.BatchNum = 1
			ret.TraceNum = 1
			err = account.Db.Model(account).Create(ret).Error
			if err != nil {
				return nil, err
			}
			return ret, err
		}
		return nil, err
	}

	return ret, nil
}

func (account *Account) IncTraceNum() error {
	account.TraceNum = (account.TraceNum + 1) % 1000000
	if account.TraceNum == 0 {
		account.TraceNum = 1
	}

	return account.Db.Model(account).Update(map[string]interface{}{"trace_num": account.TraceNum}).Error
}

func (account *Account) IncBatchNum() error {
	account.BatchNum = (account.BatchNum + 1) % 1000000
	if account.BatchNum == 0 {
		account.BatchNum = 1
	}

	return account.Db.Model(account).Update(map[string]interface{}{"batch_num": account.BatchNum}).Error
}
