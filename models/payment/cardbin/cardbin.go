package cardbin

import (
	"tpayment/models"

	"gorm.io/gorm"
)

var Dao = &CardBin{}

type CardBin struct {
	models.BaseModel
	CardNumberPrefix string `json:"card_number_prefix"`
	CardNumberLen    int    `json:"card_number_len"`
	CardNumberLuhn   bool   `json:"card_number_luhn"`
	Scheme           string `json:"scheme"`
	Type             string `json:"type"`
	Brand            string `json:"brand"`
	Prepaid          bool   `json:"prepaid"`
	CountryNumeric   string `json:"country_numeric"`
	CountryAlpha2    string `json:"country_alpha2"`
	CountryName      string `json:"country_name"`
	CountryCurrency  string `json:"country_currency"`
	CountryLatitude  string `json:"country_latitude"`
	CountryLongitude string `json:"country_longitude"`
	BankName         string `json:"bank_name"`
	BankUrl          string `json:"bank_url"`
	BankPhone        string `json:"bank_phone"`
	BankCity         string `json:"bank_city"`
}

func (CardBin) TableName() string {
	return "payment_card_bin"
}

func (k *CardBin) Create(key *CardBin) error {
	return models.DB.Model(k).Create(key).Error
}

func (k *CardBin) IsExist(cardNumberPrefix string) (bool, error) {
	ret := new(CardBin)
	err := models.DB.Model(k).Select("id").
		Where("card_number_prefix=?", cardNumberPrefix).First(ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return false, nil
		}
		return false, err
	}

	return true, nil
}
