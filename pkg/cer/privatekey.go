package cer

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func ParsePrivateKeyFromPem(priKeyData, pwd string) (crypto.PrivateKey, error) {
	var err error

	block, _ := pem.Decode([]byte(priKeyData))
	if block == nil {
		return nil, errors.New("decode fail")
	}

	if x509.IsEncryptedPEMBlock(block) {
		block.Bytes, err = x509.DecryptPEMBlock(block, []byte(pwd))
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("key type->", block.Type)
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		return x509.ParsePKCS8PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	}

	return nil, errors.New("not support key type->" + block.Type)
}
