package algorithmutils

import (
	"fmt"
	"testing"
	"tpayment/internal/encryption"
)

func TestHmac(t *testing.T) {
	ret := Hmac(encryption.BaseKey(), "50332952-bd89-4521-ad16-d43f2fdb89ca")

	fmt.Println("ret->", ret)
}
