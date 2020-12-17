package convert_utils

import (
	"fmt"
	"testing"
)

func TestASCII2ECDIC(t *testing.T) {
	fmt.Println("ASCII2ECDIC ", Bytes2HexString(ASCII2ECDIC(String2PString("123456"))))
}

func TestECDIC2ASCII(t *testing.T) {
	fmt.Println("ECDIC2ASCII ", ECDIC2ASCII(HexString2Bytes("F1F2F3F4F5F6")))
}

func TestLong2BytesEBCDIC(t *testing.T) {
	var testInt uint64 = 1234
	ECDIC := Long2BytesEBCDIC(testInt, 4)
	fmt.Println("TestLong2BytesEBCDIC ", Bytes2HexString(ECDIC))

	retInt := BytesEBCDIC2Long(ECDIC, 0, len(ECDIC))
	fmt.Println("TestBytesEBCDIC2Long ", retInt)

	if testInt != retInt {
		t.Error("testInt != retInt")
		return
	}

}
