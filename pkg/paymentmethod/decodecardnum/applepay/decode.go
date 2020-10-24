package applepay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"tpayment/pkg/algorithmutils"
)

// Rsa
func DecodeRsa(orgBean *applePayOrgBean, privateKey string, privateKeyPwd string) ([]byte, error) {
	privateKeyBean, err := parseRsaPriKey(privateKey, privateKeyPwd)
	if err != nil {
		return nil, err
	}

	wrappedKeyBytes, err := base64.StdEncoding.DecodeString(orgBean.Header.WrappedKey)
	if err != nil {
		return nil, err
	}

	wrapKeyPlain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKeyBean, wrappedKeyBytes, []byte(""))
	if err != nil {
		fmt.Println("DecryptOAEP fail->", err.Error())
		return nil, err
	}

	dataBytes, err := base64.StdEncoding.DecodeString(orgBean.Data)
	if err != nil {
		return nil, err
	}

	return algorithmutils.AESGCMDecrypt(dataBytes, wrapKeyPlain, make([]byte, 16))
}

func parseRsaPriKey(priKeyData string, privateKeyPwd string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(priKeyData))
	if x509.IsEncryptedPEMBlock(block) {
		block.Bytes, _ = x509.DecryptPEMBlock(block, []byte(privateKeyPwd))
	}

	priKeyBean, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priKeyBean.(*rsa.PrivateKey), nil
}

// ECC
//func DecodeEcc(orgBean *applePayOrgBean, privateKey, puk string, keyPwd string) ([]byte, error) {
//	privateKeyBean, err := parseEccKey(privateKey, puk, keyPwd)
//	if err != nil {
//		return nil, err
//	}
//
//	ephemeralPubKey, err := parseEccPublicKey(orgBean)
//	if err != nil {
//		return nil, err
//	}
//
//	if !elliptic.P256().IsOnCurve(ephemeralPubKey.X, ephemeralPubKey.Y) {
//		return nil, errors.New("IsOnCurve fail")
//	}
//	x, _ := elliptic.P256().ScalarMult(ephemeralPubKey.X, ephemeralPubKey.Y, privateKeyBean.D.Bytes())
//	secretBytes := x.Bytes()
//
//}
//
//func parseEccKey(privateKey, puk string, keyPwd string) (*ecdsa.PrivateKey, error) {
//	key := privateKey + "\n" + puk
//	block, _ := pem.Decode([]byte(key))
//	if x509.IsEncryptedPEMBlock(block) {
//		block.Bytes, _ = x509.DecryptPEMBlock(block, []byte(keyPwd))
//	}
//
//	priKeyBean, err := x509.ParseECPrivateKey(block.Bytes)
//
//	if err != nil {
//		return nil, nil, err
//	}
//
//	return priKeyBean, nil
//}
//
//func parseEccPublicKey(orgBean *applePayOrgBean) (*ecdsa.PublicKey, error) {
//	ephemeralPubKeyBytes := []byte("-----BEGIN PUBLIC KEY-----\n" + orgBean.Header.EphemeralPublicKey + "\n-----END PUBLIC KEY-----")
//	block, _ := pem.Decode(ephemeralPubKeyBytes)
//
//	if block == nil {
//		return nil, errors.New("parse error")
//	}
//
//	puk, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//
//	return puk.(*ecdsa.PublicKey), nil
//}
