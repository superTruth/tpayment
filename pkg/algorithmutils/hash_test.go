package algorithmutils

import (
	"fmt"
	"testing"
	"tpayment/internal/encryption"
)

func TestHmac(t *testing.T) {
	ret := Hmac(encryption.BaseKey(), "123456")

	fmt.Println("ret->", ret)
}
