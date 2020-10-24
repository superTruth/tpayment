package applepay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

//{
//"applicationPrimaryAccountNumber": "370295571160496",
//"applicationExpirationDate": "200930",
//"currencyCode": "344",
//"transactionAmount": 7500,
//"deviceManufacturerIdentifier": "030010030273",
//"paymentDataType": "3DSecure",
//"paymentData": {
//"onlinePaymentCryptogram": "IIg4zmwB1NHJNWwHBAAKoDEBhgA="
//}
//}

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

type applePayOrgBean struct {
	Data   string `json:"data"`
	Header struct {
		ApplicationData    string `json:"applicationData,omitempty"`
		EphemeralPublicKey string `json:"ephemeralPublicKey,omitempty"`
		WrappedKey         string `json:"wrappedKey,omitempty"`
		PublicKeyHash      string `json:"publicKeyHash,omitempty"`
		TransactionId      string `json:"transactionId,omitempty"`
	} `json:"header"`
	Signature string `json:"signature"`
	Version   string `json:"version"`
}

const (
	RsaEncryption = "RSA_v1"
	EccEncryption = "EC_v1"
)

func DecodeApplePay(token, privateKey string) (*ApplePayBean, error) {
	// step 1 解析base64
	//tokenBytes, err := base64.StdEncoding.DecodeString(token)
	//if err != nil {
	//	return nil, err
	//}

	applePayOrgBean := new(applePayOrgBean)
	if err := json.Unmarshal([]byte(token), applePayOrgBean); err != nil {
		return nil, err
	}

	var dataPlainByte []byte
	switch applePayOrgBean.Version {
	case RsaEncryption:
		err := ValidateRsa(applePayOrgBean)
		if err != nil {
			fmt.Println("ValidateRsa->", err.Error())
			return nil, err
		}

		priKeyFile, _ := os.Open("/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/www.fang.com.key")
		priKeyBytes, _ := ioutil.ReadAll(priKeyFile)
		dataPlainByte, err = DecodeRsa(applePayOrgBean, string(priKeyBytes), "")
		if err != nil {
			fmt.Println("DecodeRsa fail->", err)
			return nil, err
		}
	case EccEncryption:
		err := ValidateEcc(applePayOrgBean)
		if err != nil {
			fmt.Println("ValidateRsa->", err.Error())
			return nil, err
		}

	default:
		return nil, errors.New("not support encryption method " + applePayOrgBean.Version)
	}

	fmt.Println("dataPlainByte->", string(dataPlainByte))

	return nil, nil
}
