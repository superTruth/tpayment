package algorithmutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
)

func Hmac(key []byte, data []byte) string {
	hashFunc := hmac.New(sha256.New, key)
	_, _ = hashFunc.Write(data)
	return hex.EncodeToString(hashFunc.Sum([]byte("")))
}

func RandomHmacKey() []byte {
	ret := make([]byte, 16)

	for i := 0; i < len(ret); i++ {
		ret[i] = byte(rand.Int())
	}

	return ret
}

func Sha256(data []byte) string {
	ret := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(ret[:])
}
