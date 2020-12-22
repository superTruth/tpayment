package logon

import (
	"math/rand"
	"strings"
	"tpayment/pkg/utils/convert_utils"
)

// 产生随机数秘钥
func generateKekPlainHexStr() string {
	kekPlainKey := make([]byte, 16)

	for i := 0; i < len(kekPlainKey); i++ {
		kekPlainKey[i] = byte(rand.Intn(255))

		if !oddOnes(kekPlainKey[i]) {
			kekPlainKey[i] ^= 0x01
		}
	}

	return convert_utils.Bytes2HexString(kekPlainKey)
}

func oddOnes(value byte) bool {
	count := 0

	var i uint = 0
	for ; i < 8; i++ {
		if (byte(uint(1)<<i) & value) == 0 {
			count++
		}
	}
	return (count % 2) != 0
}

func generateKekPKCS1HexStr(kekPlainKeyHexStr string) string {
	//1.format TMKEncryptKeyPlaintKey to PKCS#1 form.
	sb := strings.Builder{}

	sb.WriteString("00")
	//BT
	sb.WriteString("02")
	//PS
	sb.WriteString("3667E139EF47B44773481B7636E53C2A2B767DCF28BCED78F6AF7F66DF213356FB9B97229154FC779DA69D144794A8F11B41EBA06366AADBE74BE9FD2693446A2F7E3FF1672848F245956FCB5C17C3BC73C9648358B44DA475128A9A046465")
	sb.WriteString("00")
	//D data
	sb.WriteString("30") // Fixed field
	sb.WriteString(convert_utils.Bytes2HexString(convert_utils.Long2BytesHex(uint64(1+1+len(kekPlainKeyHexStr)/2+2+8), 1)))
	sb.WriteString("04") // Fixed field
	sb.WriteString(convert_utils.Bytes2HexString(convert_utils.Long2BytesHex(uint64(len(kekPlainKeyHexStr)/2), 1)))
	sb.WriteString(kekPlainKeyHexStr)
	sb.WriteString("0408") // Fixed field
	sb.WriteString("082A89D69EC7E692")

	return sb.String()
}
