package applepay

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"tpayment/pkg/algorithmutils"
	"tpayment/pkg/cer"
	"tpayment/pkg/paymentmethod/decodecardnum/applepay/ecdh"
)

// Rsa
func decodeRsa(orgBean *applePayOrgBean, key *ConfigKey) ([]byte, error) {
	privateKeyBean, err := cer.ParsePrivateKeyFromPem(key.PrivateKey, "")
	if err != nil {
		return nil, err
	}

	wrappedKeyBytes, err := base64.StdEncoding.DecodeString(orgBean.Header.WrappedKey)
	if err != nil {
		return nil, err
	}

	wrapKeyPlain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKeyBean.(*rsa.PrivateKey), wrappedKeyBytes, []byte(""))
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

// ECC
func decodeEcc(orgBean *applePayOrgBean, key *ConfigKey) ([]byte, error) {
	wrapKeyPlain, err := eccApplePayDecode(orgBean, key)
	if err != nil {
		fmt.Println("parse key fail->", err.Error())
		return nil, err
	}

	dataBytes, err := base64.StdEncoding.DecodeString(orgBean.Data)
	if err != nil {
		return nil, err
	}

	return algorithmutils.AESGCMDecrypt(dataBytes, wrapKeyPlain, make([]byte, 16))
}

func eccApplePayDecode(orgBean *applePayOrgBean, key *ConfigKey) ([]byte, error) {
	// private key
	privateKeyBean, err := cer.ParsePrivateKeyFromPem(key.PrivateKey+"\n"+key.PublicKey, "")
	if err != nil {
		return nil, errors.New("parse private key fail->" + err.Error())
	}
	privateKey := privateKeyBean.(*ecdsa.PrivateKey)

	// public key
	publicKeyCer, err := cer.ParseCerFromPem(key.PublicKey)
	if err != nil {
		return nil, errors.New("parse public key fail->" + err.Error())
	}

	// ephemeral PubKey
	ephemeralPubKeyBytes := "-----BEGIN PUBLIC KEY-----\n" +
		orgBean.Header.EphemeralPublicKey +
		"\n-----END PUBLIC KEY-----"

	ephPublicKeyBean, err := cer.ParsePublicKeyFromPuk(ephemeralPubKeyBytes)
	if err != nil {
		return nil, errors.New("parse ephemeral public key fail->" + err.Error())
	}
	ephPublicKey := ephPublicKeyBean.(*ecdsa.PublicKey)

	// 解密ECC加密传输数据
	ec := ecdh.NewEllipticECDH(elliptic.P256())

	secretBytes, err := ec.GenerateSharedSecret(privateKey, ephPublicKey)
	if err != nil {
		return nil, errors.New("generateSharedSecret error: ->" + err.Error())
	}

	// 提取商户ID
	var merchantInfo []byte
	for _, info := range publicKeyCer.Extensions {
		id := info.Id.String()
		if id == MerchantOID {
			merchantInfo = info.Value
		}
	}
	if len(merchantInfo) == 0 {
		return nil, errors.New("can't get merchant information")
	}

	merchantName := ""
	_, err = asn1.Unmarshal(merchantInfo, &merchantName)
	if err != nil {
		return nil, errors.New("parse merchant id fail->" + err.Error())
	}

	merchantID, err := hex.DecodeString(merchantName)
	if err != nil {
		return nil, errors.New("decode merchant name hex fail->" + err.Error())
	}

	// 计算AES秘钥
	aesKey := calcAesKey(merchantID, secretBytes)

	return aesKey, nil
}

func calcAesKey(merchantID, secret []byte) []byte {
	hashMethod := sha256.New()
	_, _ = hashMethod.Write([]byte{0, 0, 0, 1})
	_, _ = hashMethod.Write(secret)
	_, _ = hashMethod.Write([]byte("\x0Did-aes256-GCM"))
	_, _ = hashMethod.Write([]byte("Apple"))
	_, _ = hashMethod.Write(merchantID)
	key := hashMethod.Sum(nil)
	return key
}
