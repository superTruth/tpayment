package basekey

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"tpayment/pkg/algorithmutils"
)

var baseKey []byte

func Init() {
	var err error
	baseKey, err = hex.DecodeString("bf76e7707b33d2360576bdb512b77260")
	if err != nil {
		panic("key error")
	}
}

func BaseKey() []byte {
	return baseKey
}

var iv = make([]byte, 16)

func EncryptData(orgData []byte) (string, error) {
	enData, err := algorithmutils.AESGCMEncrypt(orgData, baseKey, iv)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enData), nil
}

func DecryptData(orgData string) ([]byte, error) {
	orgDataBytes, err := base64.StdEncoding.DecodeString(orgData)
	if err != nil {
		return nil, errors.New("data not base64 format")
	}

	deData, err := algorithmutils.AESGCMDecrypt(orgDataBytes, baseKey, iv)
	if err != nil {
		return nil, err
	}
	return deData, nil
}

func Hash(orgData []byte) string {
	return algorithmutils.Hmac(baseKey, orgData)
}
