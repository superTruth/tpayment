package acquirer

import "tpayment/models"

var KeyDao = &Key{}

type Key struct {
	models.BaseModel
	Tag   string `gorm:"column:tag"`
	Type  string `gorm:"column:type"`
	Value string `gorm:"column:value"`
}

func (Key) TableName() string {
	return "payment_acquirer_key"
}

func (k *Key) Create(key *Key) error {
	return models.DB().Model(k).Create(key).Error
}

func (k *Key) Get(tag string) ([]*Key, error) {
	var (
		keys []*Key
		err  error
	)

	err = models.DB().Model(k).Where("tag=?",
		tag).Find(&keys).Error
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (k *Key) Delete() error {
	return models.DB().Model(k).Delete(k).Error
}
