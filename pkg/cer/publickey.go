package cer

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParsePublicKeyFromPem(publicKeyData string) (crypto.PrivateKey, error) {
	cer, err := ParseCerFromPem(publicKeyData)
	if err != nil {
		return nil, err
	}
	return cer.PublicKey, nil
}

func ParsePublicKeyFromPuk(publicKeyData string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(publicKeyData))

	if block == nil {
		return nil, errors.New("block public error")
	}

	return x509.ParsePKIXPublicKey(block.Bytes)
}

func ParseCerFromPem(publicKeyData string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(publicKeyData))
	if block == nil {
		return nil, errors.New("decode fail")
	}

	cer, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.New("parseCertificate fail->" + err.Error())
	}
	return cer, nil
}

