package encryption

import "encoding/hex"

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
