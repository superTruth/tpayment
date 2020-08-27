package algorithmutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hmac(key []byte, data string) string {
	hashFunc := hmac.New(sha256.New, key)
	_, _ = hashFunc.Write([]byte(data))
	return hex.EncodeToString(hashFunc.Sum([]byte("")))
}
