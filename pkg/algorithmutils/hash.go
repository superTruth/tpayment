package algorithmutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
)

func Hmac(key []byte, data string) string {
	hashFunc := hmac.New(sha256.New, key)
	_, _ = hashFunc.Write([]byte(data))
	return hex.EncodeToString(hashFunc.Sum([]byte("")))
}

func RandomHmacKey() []byte {
	ret := make([]byte, 16)

	for i := 0; i < len(ret); i++ {
		ret[i] = byte(rand.Int())
	}

	return ret
}
