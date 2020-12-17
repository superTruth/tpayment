package iso8583

import (
	"fmt"
	"testing"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/utils/convert_utils"
)

func TestCreateConfigFactoryFromFile(t *testing.T) {
	factory, err := CreateConfigFactory("/Users/truth/project/go/src/PaymentApp/payment-sic/iso8583config/sic_iso8583.xml")
	if err != nil {
		t.Error("CreateConfigFactory fail->", err.Error())
		return
	}

	for i, v := range factory.Configs {
		fmt.Printf("config->%#v,// F%d\n", v, i)
	}

	msg := factory.GenerateNewMessage()

	//msg.ParseMessage(convertUtils.HexString2Bytes("0210703E00810AD0829116621094688815000719000000000000010000001115523111063010110682082549034438333130313538343137323530303030303033303032353439353831303538313235303031222020202020202020202020323534393033343420202033343400179F3602016A910A0EDCC555A5FE8E433030001155543030363834313732350014220000010006004443373236434544"), 0);

	msg.SetHeaderValue("010203040506")

	msg.SetMessageType("0200")

	msg.SetFieldValue(2, "123456789")
	//msg.SetFieldVaue(3, "123456789");
	//msg.SetFieldVaue(4, "123456789");
	//field54 := "123456789";
	//msg.SetFieldVaue(54, &field54);
	//field64 := "AABBCCDD";
	msg.SetFieldValue(65, "123456")

	ret, err := msg.FormMessage()
	if err != nil {
		fmt.Println("FormMessageX err->", err)
		return
	}
	fmt.Println("FormMessageX->", len(ret), "->", convert_utils.Bytes2HexStringX(ret, ","))

	fmt.Println(msg.String())
	// parse message
	parseMsg := factory.GenerateNewMessage()
	err = parseMsg.ParseMessage(ret, 0)
	if err != nil {
		fmt.Println("parse error->", err)
		return
	}
	fmt.Println("parse message")
	fmt.Println(parseMsg.String())

}

func TestCreateConfigFactoryFromConfig(t *testing.T) {
	factory := &Factory{Configs: [129]*iso8583Define.FieldConfig{
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 4, LenLen: 0, LenType: 1, PaddingByte: 0x30, Mask: false},    // F0
		{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 16, LenLen: 0, LenType: 1, PaddingByte: 0x30, Mask: false},   // F1
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 19, LenLen: 2, LenType: 1, PaddingByte: 0x0, Mask: true},    // F2
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 6, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F3
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F4
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F5
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F6
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F7
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F8
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F9
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F10
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 6, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F11
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 6, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F12
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 4, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F13
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 4, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F14
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 4, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F15
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F16
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F17
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F18
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F19
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F20
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F21
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 3, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F22
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 3, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F23
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 3, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F24
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 2, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F25
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 2, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F26
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F27
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F28
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F29
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F30
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F31
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 11, LenLen: 2, LenType: 1, PaddingByte: 0x0, Mask: false},   // F32
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F33
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F34
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 37, LenLen: 2, LenType: 1, PaddingByte: 0x0, Mask: true},    // F35
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 104, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F36
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F37
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 6, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F38
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 2, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F39
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F40
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 8, LenLen: 0, LenType: 1, PaddingByte: 0x20, Mask: false},    // F41
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 15, LenLen: 0, LenType: 1, PaddingByte: 0x20, Mask: false},   // F42
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F43
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 25, LenLen: 2, LenType: 1, PaddingByte: 0x0, Mask: false},   // F44
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F45
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F46
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F47
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 322, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F48
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 3, LenLen: 0, LenType: 1, PaddingByte: 0x30, Mask: false},    // F49
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F50
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F51
		{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 8, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F52
		{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 16, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F53
		{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 20, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},   // F54
		{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 255, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: true},   // F55
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F56
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F57
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F58
		{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 255, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F59
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 17, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},   // F60
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 29, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},   // F61
		{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 99, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},   // F62
		{FieldValueType: 0, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 163, LenLen: 3, LenType: 1, PaddingByte: 0x20, Mask: false}, // F63
		{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 8, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},     // F64
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F65
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F66
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F67
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F68
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F69
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F70
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F71
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F72
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F73
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F74
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F75
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F76
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F77
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F78
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F79
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F80
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F81
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F82
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F83
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F84
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F85
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F86
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F87
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F88
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F89
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F90
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F91
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F92
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F93
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F94
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F95
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0xf1, Mask: false},   // F96
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0xf1, Mask: false},   // F97
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F98
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F99
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F100
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F101
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F102
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F103
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F104
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F105
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F106
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F107
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F108
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F109
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F110
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F111
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F112
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F113
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F114
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F115
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F116
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F117
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F118
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F119
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F120
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F121
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F122
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F123
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F124
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F125
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F126
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F127
		{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F128
	}}

	msg := factory.GenerateNewMessage()

	msg.SetHeaderValue("010203040506")

	msg.SetMessageType("0200")

	msg.SetFieldValue(2, "123456789")
	msg.SetFieldValue(65, "123456")

	ret, err := msg.FormMessage()
	if err != nil {
		fmt.Println("FormMessageX err->", err)
		return
	}
	fmt.Println("FormMessageX->", len(ret), "->", convert_utils.Bytes2HexStringX(ret, ","))

	fmt.Println(msg.String())
	// parse message
	parseMsg := factory.GenerateNewMessage()
	err = parseMsg.ParseMessage(ret, 0)
	if err != nil {
		fmt.Println("parse error->", err)
		return
	}
	fmt.Println("parse message")
	fmt.Println(parseMsg.String())

}
