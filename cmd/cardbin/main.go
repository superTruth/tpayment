package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/payment/cardbin"
	"tpayment/pkg/paymentmethod/decodecardnum/creditcard"
)

func main() {
	conf.InitConfigData()
	models.InitDB()

	for {
		cardBrand := rand.Int() % len(creditcard.Rules)
		cardNumberRange := rand.Int() % len(creditcard.Rules[cardBrand].CardNumPreFix)
		cardRange := creditcard.Rules[cardBrand].CardNumPreFix[cardNumberRange]
		cardBin := rand.Int()%(cardRange.End-cardRange.Start) + cardRange.Start
		ok, err := cardbin.Dao.IsExist(fmt.Sprintf("%06d", cardBin))
		if err != nil || ok {
			continue
		}
		err = exe(cardBin)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * 6)
	}
	//
	//for i := 1; ; {
	//	fmt.Println("===========", creditcard.Rules[i].CardBrand, "===========")
	//	for j := 0; j < len(creditcard.Rules[i].CardNumPreFix); j++ {
	//		for k := creditcard.Rules[i].CardNumPreFix[j].Start; k < creditcard.Rules[i].CardNumPreFix[j].End; k++ {
	//			ok, err := cardbin.Dao.IsExist(fmt.Sprintf("%06d", k))
	//			if err != nil || ok {
	//				continue
	//			}
	//
	//			err = exe(k)
	//			if err != nil {
	//				fmt.Println(err.Error())
	//			}
	//			time.Sleep(time.Second * 6)
	//		}
	//	}
	//
	//	i = (i + 1) % len(creditcard.Rules)
	//	if i == len(creditcard.Rules) {
	//		i = 1
	//	}
	//}
}

type BinBean struct {
	Scheme  string `json:"scheme"`
	Type    string `json:"type"`
	Brand   string `json:"brand"`
	PrePaid bool   `json:"prepaid"`

	Number  *Number  `json:"number"`
	Country *Country `json:"country"`
	Bank    *Bank    `json:"bank"`
}

type Number struct {
	Length int  `json:"length"`
	Luhn   bool `json:"luhn"`
}

type Country struct {
	Numeric  string `json:"numeric"`
	Alpha2   string `json:"alpha_2"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}
type Bank struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Phone string `json:"phone"`
	City  string `json:"city"`
}

// 从网络获取数据
func GetCardBin(prefix int) (*BinBean, error) {
	url := "https://lookup.binlist.net/" + fmt.Sprintf("%06d", prefix)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return &BinBean{
			Scheme:  "unknown",
			Type:    "unknown",
			Brand:   "unknown",
			PrePaid: false,
			Number:  nil,
			Country: nil,
			Bank:    nil,
		}, nil
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("body->", string(body))
		return nil, errors.New("StatusCode->" + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := new(BinBean)
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// 数据转换
func change2DB(prefix int, bin *BinBean) *cardbin.CardBin {
	binDB := &cardbin.CardBin{
		CardNumberPrefix: fmt.Sprintf("%06d", prefix),
		Scheme:           bin.Scheme,
		Type:             bin.Type,
		Brand:            bin.Brand,
		Prepaid:          bin.PrePaid,
	}

	if bin.Number != nil {
		binDB.CardNumberLen = bin.Number.Length
		binDB.CardNumberLuhn = bin.Number.Luhn
	}

	if bin.Country != nil {
		binDB.CountryNumeric = bin.Country.Numeric
		binDB.CountryAlpha2 = bin.Country.Alpha2
		binDB.CountryName = bin.Country.Name
		binDB.CountryCurrency = bin.Country.Currency
	}

	if bin.Bank != nil {
		binDB.BankName = bin.Bank.Name
		binDB.BankUrl = bin.Bank.Url
		binDB.BankPhone = bin.Bank.Phone
		binDB.BankCity = bin.Bank.City
	}

	return binDB
}

// 执行操作
func exe(prefix int) error {
	fmt.Println("exe->", prefix)
	var dbBin *cardbin.CardBin
	cardInfo, err := GetCardBin(prefix)
	if err != nil {
		return errors.New("GetCardBin error->" + err.Error())
	} else {
		dbBin = change2DB(prefix, cardInfo)
	}
	err = cardbin.Dao.Create(dbBin)
	if err != nil {
		return errors.New("Create error->" + err.Error())
	}
	return nil
}
