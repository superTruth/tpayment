package applepay

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/zhulingbiezhi/pkcs7"
)

func ValidateRsa(orgBean *applePayOrgBean) error {
	fmt.Println("Signature->", orgBean.Signature)
	signData, err := base64.StdEncoding.DecodeString(orgBean.Signature)
	if err != nil {
		fmt.Println("DecodeString fail->", err.Error())
		return err
	}

	p7, err := GenerateValidatePkcs7(signData)
	if err != nil {
		fmt.Println("GenerateValidatePkcs7 fail->", err.Error())
		return err
	}

	//
	wrappedKeyBytes, err := base64.StdEncoding.DecodeString(orgBean.Header.WrappedKey)
	if err != nil {
		fmt.Println("DecodeString2 fail->", err.Error())
		return err
	}

	dataBytes, err := base64.StdEncoding.DecodeString(orgBean.Data)
	if err != nil {
		return err
	}

	transactionIdBytes, err := hex.DecodeString(orgBean.Header.TransactionId)
	if err != nil {
		return err
	}

	applicationDataBytes, err := hex.DecodeString(orgBean.Header.ApplicationData)
	if err != nil {
		return err
	}

	sb := bytes.Buffer{}
	sb.Write(wrappedKeyBytes)
	sb.Write(dataBytes)
	sb.Write(transactionIdBytes)
	sb.Write(applicationDataBytes)

	p7.Content = sb.Bytes()

	return p7.Verify(x509.SHA256WithRSA)
}

func ValidateEcc(orgBean *applePayOrgBean) error {
	signData, err := base64.StdEncoding.DecodeString(orgBean.Signature)
	if err != nil {
		return err
	}

	p7, err := GenerateValidatePkcs7(signData)
	if err != nil {
		return err
	}

	//
	ephemeralPubKeyBytes := []byte("-----BEGIN PUBLIC KEY-----\n" + orgBean.Header.EphemeralPublicKey + "\n-----END PUBLIC KEY-----")
	pukByte, _ := pem.Decode(ephemeralPubKeyBytes)
	if pukByte == nil {
		return errors.New("decode EphemeralPublicKey fail")
	}

	dataBytes, err := base64.StdEncoding.DecodeString(orgBean.Data)
	if err != nil {
		return err
	}

	transactionIdBytes, err := hex.DecodeString(orgBean.Header.TransactionId)
	if err != nil {
		return err
	}

	applicationDataBytes, err := hex.DecodeString(orgBean.Header.ApplicationData)
	if err != nil {
		return err
	}

	sb := bytes.Buffer{}
	sb.Write(ephemeralPubKeyBytes)
	sb.Write(dataBytes)
	sb.Write(transactionIdBytes)
	sb.Write(applicationDataBytes)

	p7.Content = sb.Bytes()

	return p7.Verify(x509.SHA256WithRSA)
}

func GenerateValidatePkcs7(signData []byte) (*pkcs7.PKCS7, error) {
	p7, err := pkcs7.Parse(signData)
	if err != nil {
		return p7, err
	}
	validateCount := 0
	for _, cer := range p7.Certificates {
		for _, info := range cer.Extensions {
			idStr := info.Id.String()
			if idStr == LeafOID || idStr == CAOID {
				validateCount++
				break
			}
		}
		if validateCount == 2 {
			break
		}
	}
	if validateCount < 2 {
		return p7, errors.New("certification not correct")
	}

	return p7, nil
}
