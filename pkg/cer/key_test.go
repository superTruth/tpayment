package cer

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParsePrivateKeyFromPem(t *testing.T) {

	// rsa
	const privateFilePath = "/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/www.fang.com.key"
	priKeyFile, _ := os.Open(privateFilePath)
	priKeyBytes, _ := ioutil.ReadAll(priKeyFile)

	_, err := ParsePrivateKeyFromPem(string(priKeyBytes), "")
	if err != nil {
		t.Error("rsa fail->", err.Error())
		return
	}


	// ecc
	const eccPukFilePath = "/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/apple_pay_ecc.crt"
	const eccPrivateFilePath = "/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/ecc.key"

	eccPukFile, _ := os.Open(eccPukFilePath)
	eccPukBytes, _ := ioutil.ReadAll(eccPukFile)

	eccPriKeyFile, _ := os.Open(eccPrivateFilePath)
	eccPriKeyBytes, _ := ioutil.ReadAll(eccPriKeyFile)

	_, err = ParsePrivateKeyFromPem(string(eccPriKeyBytes) + "\n" + string(eccPukBytes), "")
	if err != nil {
		t.Error("ecc fail->", err.Error())
		return
	}
}

func TestParsePublicKeyFromPem(t *testing.T) {
	// rsa
	const pubKeyFilePath = "/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/apple_pay.crt"
	pukFile, _ := os.Open(pubKeyFilePath)
	pukBytes, _ := ioutil.ReadAll(pukFile)

	//pukStr := "-----BEGIN CERTIFICATE-----\n" +
	//	base64.StdEncoding.EncodeToString(pukBytes) +
	//	"\n-----END CERTIFICATE-----"
	//
	//fmt.Println("ecc puk->", pukStr)
	//
	//block, _ := pem.Decode([]byte(pukStr))
	//if block == nil {
	//	t.Error("decode fail")
	//	return
	//}
	//
	//
	//_,err := x509.ParseCertificate(pukBytes)

	_, err := ParsePublicKeyFromPem(string(pukBytes))
	if err != nil {
		t.Error("rsa fail->", err.Error())
		return
	}


	// ecc
	const eccPukFilePath = "/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/apple_pay_ecc.crt"

	eccPukFile, _ := os.Open(eccPukFilePath)
	eccPukBytes, _ := ioutil.ReadAll(eccPukFile)

	_, err = ParsePublicKeyFromPem(string(eccPukBytes))
	if err != nil {
		t.Error("ecc fail->", err.Error())
		return
	}
}
