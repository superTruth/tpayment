package algorithmutils

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/format_utils"
	"tpayment/pkg/utils/mix_utils"
)

// DES ECB加密
func EncryptDesECB(data []byte, key []byte) ([]byte, error) {

	if data == nil || key == nil {
		return nil, errors.New("encrypt parameters error")
	}

	keyBlock, err := generateCipherBlock(key)
	if err != nil {
		return nil, err
	}

	formatedDatas, err := formatData(data)
	if err != nil {
		return nil, err
	}

	tmpBlock := make([]byte, 8)
	encryptBlock := make([]byte, 8)

	for i := 0; i < len(formatedDatas)/8; i++ {
		mix_utils.BytesArrayCopy(formatedDatas, 8*i, tmpBlock, 0, 8)
		keyBlock.Encrypt(encryptBlock, tmpBlock)
		mix_utils.BytesArrayCopy(encryptBlock, 0, formatedDatas, 8*i, 8)
	}

	return formatedDatas, nil
}

// DES ECB解密
func DecryptDesECB(data []byte, key []byte) ([]byte, error) {
	if data == nil || key == nil {
		return nil, errors.New("encrypt parameters error")
	}

	keyBlock, err := generateCipherBlock(key)
	if err != nil {
		return nil, err
	}

	formatedDatas, err := formatData(data)
	if err != nil {
		return nil, err
	}

	tmpBlock := make([]byte, 8)
	decryptBlock := make([]byte, 8)
	for i := 0; i < len(formatedDatas)/8; i++ {
		mix_utils.BytesArrayCopy(formatedDatas, 8*i, tmpBlock, 0, 8)

		keyBlock.Decrypt(decryptBlock, tmpBlock)

		mix_utils.BytesArrayCopy(decryptBlock, 0, formatedDatas, 8*i, 8)
	}

	return formatedDatas, nil
}

// DES CBC加密
func EncryptDesCBC(data []byte, key []byte, iv []byte) ([]byte, error) {

	if (iv == nil) || (len(iv) != 8) {
		return nil, errors.New("iv len must 8")
	}

	keyBlock, err := generateCipherBlock(key)
	if err != nil {
		return nil, err
	}

	formatedDatas, err := formatData(data)
	if err != nil {
		return nil, err
	}

	cipher.NewCBCEncrypter(keyBlock, iv).CryptBlocks(formatedDatas, formatedDatas)

	return formatedDatas, nil
}

// DES ECB解密
func DecryptDesCBC(data []byte, key []byte, iv []byte) ([]byte, error) {

	if (iv == nil) || (len(iv) != 8) {
		return nil, errors.New("iv len must 8")
	}

	keyBlock, err := generateCipherBlock(key)
	if err != nil {
		return nil, err
	}

	formatedDatas, err := formatData(data)
	if err != nil {
		return nil, err
	}

	cipher.NewCBCDecrypter(keyBlock, iv).CryptBlocks(formatedDatas, formatedDatas)

	return formatedDatas, nil
}

// 生成key
func generateCipherBlock(key []byte) (cipher.Block, error) {
	if (key == nil) || ((len(key) != 8) && (len(key) != 16) && (len(key) != 24)) {
		return nil, errors.New("key len error")
	}

	var ret cipher.Block
	var err error

	switch len(key) {
	case 8:
		ret, err = des.NewCipher(key)
		break
	case 16:
		newKey := make([]byte, 24)
		mix_utils.BytesArrayCopy(key, 0, newKey, 0, 16)
		mix_utils.BytesArrayCopy(key, 0, newKey, 16, 8)
		ret, err = des.NewTripleDESCipher(newKey)
		break
	case 24:
		ret, err = des.NewTripleDESCipher(key)
		break
	}

	return ret, err
}

// 格式化原始数据
func formatData(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data can't nil")
	}

	formatedData := make([]byte, (len(data)+7)/8*8)
	mix_utils.BytesArrayCopy(data, 0, formatedData, 0, len(data))
	return formatedData, nil
}

// RSA加密
func RsaPublicEncryption(data []byte, modulus *string, exponent *string) ([]byte, error) {

	modulusObject, ok := new(big.Int).SetString(*modulus, 16)
	if !ok {
		return nil, errors.New("bad modulus")
	}

	exponentObject, err := strconv.ParseInt(*exponent, 16, 64)
	if err != nil {
		return nil, errors.New("bad exponent")
	}

	pk := rsa.PublicKey{
		N: modulusObject,
		E: int(exponentObject),
	}

	fmt.Println("key len->", pk.Size())

	random := rand.Reader

	ret, err := rsa.EncryptPKCS1v15(random, &pk, data)

	return ret, err
}

func RsaPublicEncryptionX(data []byte, modulus string, exponent string) ([]byte, error) {

	modulusObject, _ := new(big.Int).SetString(modulus, 16)
	exponentObject, _ := new(big.Int).SetString(exponent, 16)
	messageObject := new(big.Int).SetBytes(data)

	ret := new(big.Int).Exp(messageObject, exponentObject, modulusObject)

	return ret.Bytes(), nil
}

// 校验kcv
func CheckKCV(value []byte, kcv []byte) bool {
	if (value == nil) || (kcv == nil) || (len(value)%8 != 0) || (len(kcv) != 4) {
		return false
	}

	calcKcv := CalcKcv(value)
	//zeroBytes := convert_utils.HexString2Bytes("0000000000000000");
	//calcKcv, _ := EncryptDesECB(zeroBytes, value);
	//calcKcv = mix_utils.BytesArrayCopyArrange(calcKcv, 0, 4);
	return mix_utils.Compare(kcv, calcKcv)
}

func CalcKcv(value []byte) []byte {
	if (value == nil) || (len(value)%8 != 0) {
		return nil
	}
	zeroBytes := convert_utils.HexString2Bytes("0000000000000000")
	calcKcv, _ := EncryptDesECB(zeroBytes, value)
	calcKcv = mix_utils.BytesArrayCopyArrange(calcKcv, 0, 4)

	return calcKcv
}

// 加密PIN
func EncryptPIN(plainPin string, PAN string, PINKey []byte) ([]byte, error) {
	if len(PAN) < 13 || len(PAN) > 19 {
		return nil, errors.New("PAN format error")
	}
	PANBlock := "0000" + mix_utils.SubString(PAN, len(PAN)-13, len(PAN)-1)
	if len(plainPin) == 0 {
		return []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, nil
	}

	PINBlock := fmt.Sprintf("%02d", len(plainPin)) + plainPin
	PINBlock = format_utils.AppendString(PINBlock, 16, false, 'F')
	PANBlockByte := convert_utils.HexString2Bytes(PANBlock)
	PINBlockByte := convert_utils.HexString2Bytes(PINBlock)
	ret, err := mix_utils.Xor(PANBlockByte, PINBlockByte)
	if err != nil {
		return nil, err
	}

	return EncryptDesECB(ret, PINKey)
}

// 计算MAC CUP方式
func CalcCUPMac(value []byte, MACKey []byte) ([]byte, error) {
	if len(value) == 0 || (len(MACKey)%8 != 0) {
		return nil, errors.New("encrypt parameters error")
	}

	var tmpBlockByte []byte
	realMacKey := MACKey //mix_utils.BytesArrayCopyArrange(MACKey, 0, 8);

	formatValue, err := formatData(value)
	if err != nil {
		return nil, err
	}

	// step 1 循环进行异或计算
	preBlock := make([]byte, 8)
	for i := 0; i < len(formatValue)/8; i++ {
		mix_utils.BytesArrayCopy(formatValue, i*8, preBlock, 0, 8)
		if tmpBlockByte == nil {
			tmpBlockByte = make([]byte, 8)
			mix_utils.BytesArrayCopy(formatValue, i*8, tmpBlockByte, 0, 8)
			continue
		} else {
			mix_utils.BytesArrayCopy(formatValue, i*8, preBlock, 0, 8)
		}

		tmpBlockByte, _ = mix_utils.Xor(tmpBlockByte, preBlock)
	}

	// step 2 结果换算成16字节
	tmpBlockByte = []byte(convert_utils.Bytes2HexString(tmpBlockByte))

	// step 3 提取前面8个字节，进行一次加密
	blockByte1 := mix_utils.BytesArrayCopyArrange(tmpBlockByte, 0, 8)
	blockByte1, _ = EncryptDesECB(blockByte1, realMacKey)

	// step 4 提取后面8个字节和结果异或
	blockByte2 := mix_utils.BytesArrayCopyArrange(tmpBlockByte, 8, 16)
	tmpBlockByte, _ = mix_utils.Xor(blockByte1, blockByte2)

	// step 5 结果进行加密
	tmpBlockByte, _ = EncryptDesECB(tmpBlockByte, realMacKey)
	tmpBlockByte = []byte(convert_utils.Bytes2HexString(tmpBlockByte))
	return mix_utils.BytesArrayCopyArrange(tmpBlockByte, 0, 8), nil
}

// 计算MAC ANSI9.9方式
func CalcMacANSI99(value []byte, MACKey []byte, iv []byte) ([]byte, error) {
	if value == nil || MACKey == nil || (len(MACKey)%8 != 0 || iv == nil || len(iv) != 8) {
		return nil, errors.New("encrypt parameters error")
	}

	var tmpBlockByte []byte

	formattedValue, err := formatData(value)
	if err != nil {
		return nil, err
	}

	// step 1 循环进行异或和DES/3DES计算
	preBlock := make([]byte, 8)
	for i := 0; i < len(formattedValue)/8; i++ {
		mix_utils.BytesArrayCopy(formattedValue, i*8, preBlock, 0, 8)
		if tmpBlockByte == nil {
			tmpBlockByte = make([]byte, 8)
			mix_utils.BytesArrayCopy(formattedValue, i*8, tmpBlockByte, 0, 8)
			tmpBlockByte, _ = mix_utils.Xor(tmpBlockByte, iv)
			tmpBlockByte, err = EncryptDesECB(tmpBlockByte, MACKey)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			mix_utils.BytesArrayCopy(formattedValue, i*8, preBlock, 0, 8)
		}
		//fmt.Println("i = ", i, ", tempBlock = ", convert_utils.Bytes2HexString(tmpBlockByte), ", preBlock = ", convert_utils.Bytes2HexString(preBlock))
		tmpBlockByte, _ = mix_utils.Xor(tmpBlockByte, preBlock)
		tmpBlockByte, err = EncryptDesECB(tmpBlockByte, MACKey)
		if err != nil {
			return nil, err
		}
	}

	return tmpBlockByte, nil
}

// 计算MAC ANSI X9.19
func CalcMacANSI919(value []byte, MACKey []byte, iv []byte) ([]byte, error) {
	if value == nil || MACKey == nil || (len(MACKey) != 16 || iv == nil || len(iv) != 8) {
		return nil, errors.New("encrypt parameters error")
	}

	var keyLeft = make([]byte, 8)
	var keyRight = make([]byte, 8)

	mix_utils.BytesArrayCopy(MACKey, 0, keyLeft, 0, 8)
	mix_utils.BytesArrayCopy(MACKey, 8, keyRight, 0, 8)

	// 用MAC密钥左半部做 X9.9算法
	result99, err := CalcMacANSI99(value, keyLeft, iv)
	if err != nil {
		return nil, err
	}

	//用MAC密钥右半部解密result 99
	resultTemp, err := DecryptDesECB(result99, keyRight)
	if err != nil {
		return nil, err
	}

	//用MAC密钥左半部加密resultTemp，结果为最终结果
	return EncryptDesECB(resultTemp, keyLeft)

}

// 加密TK2
func EncryptTK(plainTK string, TDK []byte) (string, error) {
	formatedTk := strings.Replace(plainTK, "=", "D", -1)
	if len(formatedTk)%2 != 0 { // 不是2的倍数，则后边填充F
		formatedTk = formatedTk + "F"
	}
	tkBytes := convert_utils.HexString2Bytes(formatedTk)
	if len(tkBytes) < 10 {
		return "", errors.New("tk too short")
	}

	// 截取8个字节，进行加密
	tdb := mix_utils.BytesArrayCopyArrange(tkBytes, len(tkBytes)-9, len(tkBytes)-1)
	tdbEn, _ := EncryptDesECB(tdb, TDK)

	// 将结果放回原处
	mix_utils.BytesArrayCopy(tdbEn, 0, tkBytes, len(tkBytes)-9, len(tdbEn))
	ret := convert_utils.Bytes2HexString(tkBytes)
	// 截取掉最后一位F
	if len(plainTK)%2 != 0 {
		return mix_utils.SubString(ret, 0, len(ret)-1), nil
	} else {
		return ret, nil
	}
}
