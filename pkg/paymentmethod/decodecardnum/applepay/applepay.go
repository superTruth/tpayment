package applepay

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ApplePayBean struct {
	ApplicationPrimaryAccountNumber string `json:"applicationPrimaryAccountNumber"`
	ApplicationExpirationDate       string `json:"applicationExpirationDate"` // YYMMDD
	CurrencyCode                    string `json:"currencyCode"`
	TransactionAmount               string `json:"transactionAmount"`
	CardholderName                  string `json:"cardholderName"`
	DeviceManufacturerIdentifier    string `json:"deviceManufacturerIdentifier"`
	PaymentDataType                 string `json:"paymentDataType"`
}

type PaymentData struct {
	OnlinePaymentCryptogram string `json:"onlinePaymentCryptogram"`
	EciIndicator            string `json:"eciIndicator"`
	EmvData                 string `json:"emvData"`
	EncryptedPINData        string `json:"encryptedPINData"`
}

type headerBean struct {
	ApplicationData    string `json:"applicationData,omitempty"`
	EphemeralPublicKey string `json:"ephemeralPublicKey,omitempty"`
	WrappedKey         string `json:"wrappedKey,omitempty"`
	PublicKeyHash      string `json:"publicKeyHash,omitempty"`
	TransactionId      string `json:"transactionId,omitempty"`
}
type applePayOrgBean struct {
	Data      string      `json:"data"`
	Header    *headerBean `json:"header"`
	Signature string      `json:"signature"`
	Version   string      `json:"version"`
}

const (
	RsaEncryption = "RSA_v1"
	EccEncryption = "EC_v1"
)

func covertApplePayToken(token string) (*applePayOrgBean, error) {
	ret := new(applePayOrgBean)
	if err := json.Unmarshal([]byte(token), ret); err != nil {
		return nil, err
	}

	if ret.Data == "" || ret.Header == nil || ret.Signature == "" || ret.Version == "" ||
		ret.Header.PublicKeyHash == "" {
		return nil, errors.New("token is invalidation")
	}

	return ret, nil
}

// 获取使用的Key Hash
func GetApplePayKeyHash(token string) (string, error) {
	applePayOrgBean, err := covertApplePayToken(token)
	if err != nil {
		return "", err
	}

	return applePayOrgBean.Header.PublicKeyHash, err
}

type ConfigKey struct {
	PublicKey  string
	PrivateKey string
}

func DecodeApplePay(token string, key *ConfigKey) (*ApplePayBean, error) {
	applePayOrgBean, err := covertApplePayToken(token)
	if err != nil {
		return nil, err
	}

	var dataPlainByte []byte
	switch applePayOrgBean.Version {
	case RsaEncryption:
		err := validateSignature(applePayOrgBean, validateRsa)
		if err != nil {
			fmt.Println("validateRsa->", err.Error())
			return nil, err
		}

		dataPlainByte, err = DecodeRsa(applePayOrgBean, key)
		if err != nil {
			fmt.Println("DecodeRsa fail->", err)
			return nil, err
		}
	case EccEncryption:
		fmt.Println("ecc")
		err := validateSignature(applePayOrgBean, validateEcc)
		if err != nil {
			fmt.Println("validateRsa->", err.Error())
			return nil, err
		}
		dataPlainByte, err = DecodeEcc(applePayOrgBean, key)
		if err != nil {
			fmt.Println("DecodeEcc fail->", err)
			return nil, err
		}

	default:
		return nil, errors.New("not support encryption method " + applePayOrgBean.Version)
	}

	fmt.Println("dataPlainByte->", string(dataPlainByte))

	return nil, nil
}
