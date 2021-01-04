package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"tpayment/internal/acquirer_impl"
)

type Config struct {
	AcquirerConfig *AcquirerConfig
	MerchantConfig *MerchantConfig
	TerminalConfig *TerminalConfig
}

type AcquirerConfig struct {
	TPDU string `json:"TPDU,omitempty"`
	NII  string `json:"NII,omitempty"`
	MID  string `json:"MID,omitempty"`
	TID  string `json:"TID,omitempty"`
	URL  string `json:"URL,omitempty"`
}

type MerchantConfig struct {
}

type TerminalConfig struct {
}

func ParseConfig(req *acquirer_impl.SaleRequest) (*Config, error) {
	var err error
	ret := new(Config)

	if req.TxqReq.PaymentProcessRule == nil {
		return nil, errors.New("no merchant process rule")
	}

	if req.TxqReq.PaymentProcessRule.MerchantAccount == nil {
		return nil, errors.New("no merchant account")
	}

	if req.TxqReq.PaymentProcessRule.MerchantAccount.Acquirer == nil {
		return nil, errors.New("no acquirer")
	}

	ret.AcquirerConfig = new(AcquirerConfig)
	err = json.Unmarshal([]byte(req.TxqReq.PaymentProcessRule.MerchantAccount.Acquirer.Addition), ret.AcquirerConfig)
	if err != nil {
		return nil, fmt.Errorf("parse acquirer addition fail->%v", err)
	}

	return ret, nil
}
