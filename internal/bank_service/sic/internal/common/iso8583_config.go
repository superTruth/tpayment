package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/pkg/iso8583"
	"tpayment/pkg/iso8583/iso8583Define"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils/format_utils"
)

var iso8583Factory = &iso8583.Factory{Configs: [129]*iso8583Define.FieldConfig{
	{FieldValueType: 1, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 4, LenLen: 0, LenType: 1, PaddingByte: 0x30, Mask: false},    // F0
	{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: true, ValueLen: 8, LenLen: 0, LenType: 1, PaddingByte: 0x30, Mask: false},    // F1
	{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 19, LenLen: 2, LenType: 1, PaddingByte: 0x0, Mask: false},   // F2
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
	{FieldValueType: 2, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 37, LenLen: 2, LenType: 1, PaddingByte: 0x0, Mask: false},   // F35
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
	{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 255, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F55
	{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F56
	{FieldValueType: 0, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 100, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F57
	{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F58
	{FieldValueType: 0, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 255, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F59
	{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 999, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F60
	{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 999, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F61
	{FieldValueType: 2, FieldAlignType: 0, IsValueLenFix: false, ValueLen: 999, LenLen: 3, LenType: 1, PaddingByte: 0x0, Mask: false},  // F62
	{FieldValueType: 0, FieldAlignType: 1, IsValueLenFix: false, ValueLen: 999, LenLen: 3, LenType: 1, PaddingByte: 0x20, Mask: false}, // F63
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
	{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F96
	{FieldValueType: 1, FieldAlignType: 1, IsValueLenFix: true, ValueLen: 12, LenLen: 0, LenType: 1, PaddingByte: 0x0, Mask: false},    // F97
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

func SetCommonDataToMsg(req *acquirer_impl.SaleRequest, config *Config, fieldSelect []byte) (*iso8583.Message, error) {
	iso8583Msg := iso8583Factory.GenerateNewMessage()

	_ = iso8583Msg.SetHeaderValue(config.AcquirerConfig.TPDU)
	for i := 0; i < len(fieldSelect); i++ {
		switch fieldSelect[i] {
		case 2:
			if req.TxqReq.RealEntryType == conf.Contact ||
				req.TxqReq.RealEntryType == conf.ContactLess {
				_ = iso8583Msg.SetFieldValue(2, req.TxqReq.CreditCardBean.CardNumber)
			}
		case 4:
			amount, err := format_utils.FloatAmount2Int(req.TxqReq.Amount, 2)
			if err != nil {
				return nil, errors.New("amount format error->" + err.Error())
			}
			_ = iso8583Msg.SetFieldValue(4, amount)
		case 11:
			_ = iso8583Msg.SetFieldValue(11, strconv.FormatUint(req.TxqReq.CreditCardBean.TraceNum, 10))
		case 14:
			_ = iso8583Msg.SetFieldValue(14,
				req.TxqReq.CreditCardBean.CardExpYear+req.TxqReq.CreditCardBean.CardExpMonth)
		case 22:
			entryMode := strings.Builder{}
			switch req.TxqReq.RealEntryType {
			case conf.ManualInput:
				entryMode.WriteString("01")
			case conf.Swipe:
				entryMode.WriteString("02")
			case conf.Contact:
				entryMode.WriteString("05")
			case conf.ContactLess:
				entryMode.WriteString("07")
			case conf.ConsumerPresentQR:
				entryMode.WriteString("03")
			default:
				if req.TxqReq.CreditCardBean.IsMsdCard {
					entryMode.WriteString("91")
				} else {
					return nil, errors.New("not support entry mode->" + req.TxqReq.RealEntryType)
				}
			}
			if req.TxqReq.CreditCardBean.PIN == "" {
				entryMode.WriteString("2")
			} else {
				entryMode.WriteString("1")
			}
			_ = iso8583Msg.SetFieldValue(22, entryMode.String())
		case 23:
			if req.TxqReq.CreditCardBean.CardSn != "" {
				_ = iso8583Msg.SetFieldValue(23, fmt.Sprintf("%02s", req.TxqReq.CreditCardBean.CardSn))
			}
		case 26:
			if req.TxqReq.CreditCardBean.PIN != "" {
				_ = iso8583Msg.SetFieldValue(26, "12")
			}
		case 35:
			_ = iso8583Msg.SetFieldValue(35, req.TxqReq.CreditCardBean.CardTrack2)
		case 37:
			if req.TxqReq.OrgRecord != nil {
				_ = iso8583Msg.SetFieldValue(37, req.TxqReq.OrgRecord.AcquirerRRN)
			}
		case 38:
			if req.TxqReq.OrgRecord != nil {
				_ = iso8583Msg.SetFieldValue(38, req.TxqReq.OrgRecord.AcquirerAuthCode)
			}
		case 41:
			_ = iso8583Msg.SetFieldValue(41, req.TxqReq.PaymentProcessRule.MerchantAccount.Terminal.TID)
		case 42:
			_ = iso8583Msg.SetFieldValue(42, req.TxqReq.PaymentProcessRule.MerchantAccount.MID)
		case 49:
			_ = iso8583Msg.SetFieldValue(49, conf.CurrencyCode[req.TxqReq.Currency])
		case 52:
			_ = iso8583Msg.SetFieldValue(52, req.TxqReq.CreditCardBean.PIN)
		case 53:
			if req.TxqReq.CreditCardBean.PIN != "" || req.TxqReq.CreditCardBean.CardTrack2 != "" {
				f53VauleSb := strings.Builder{}
				if req.TxqReq.CreditCardBean.PIN != "" {
					f53VauleSb.WriteString("2")
				} else {
					f53VauleSb.WriteString("0")
				}
				f53VauleSb.WriteString("61")
				f53VauleSb.WriteString("0000000000000")
				_ = iso8583Msg.SetFieldValue(53, f53VauleSb.String())
			}
		case 55:
			_ = iso8583Msg.SetFieldValue(55, req.TxqReq.CreditCardBean.IccRequest)
		case 64:
			_ = iso8583Msg.SetFieldValue(64, "0000000000000000")
		}
	}

	return iso8583Msg, nil
}

func GetCommonDataFromMsg(srcData []byte, offset int) (*acquirer_impl.SaleResponse, *iso8583.Message, error) {
	logger := tlog.GetGoroutineLogger()

	resp := new(acquirer_impl.SaleResponse)
	resp.TxnResp = new(api_define.TxnResp)
	resp.TxnResp.CreditCardBean = new(api_define.CreditCardBean)
	var err error

	msg := iso8583Factory.GenerateNewMessage()
	if err = msg.ParseMessage(srcData, offset); err != nil {
		return nil, nil, errors.New("parse iso8583 fail->" + err.Error())
	}

	logger.Info("parse msg->", msg.String())

	resp.TxnResp.AcquirerRRN, _ = msg.GetFieldValue(37)
	resp.TxnResp.CreditCardBean.AuthCode, _ = msg.GetFieldValue(38)
	resp.TxnResp.CreditCardBean.ResponseCode, _ = msg.GetFieldValue(39)
	resp.TxnResp.CreditCardBean.IccResponse, _ = msg.GetFieldValue(55)

	hostTime, _ := msg.GetFieldValue(12)
	hostDate, _ := msg.GetFieldValue(13)
	if hostTime != "" && hostDate != "" {
		resp.TxnResp.DateTime, err = formatAcquirerTime(hostDate, hostTime)
		if err != nil {
			return resp, nil, errors.New("date time format error->d:" + hostDate + ", t:" + hostTime)
		}
	}

	// 银联专属UPI, 保存在RFU1里面
	if len(resp.TxnResp.AcquirerRRN) > 6 {
		// 取最后6位作为结果
		resp.TxnResp.AdditionData = new(api_define.AdditionData)
		resp.TxnResp.AdditionData.CupTraceNum = resp.TxnResp.AcquirerRRN[len(resp.TxnResp.AcquirerRRN)-6:]
	}

	return resp, msg, nil
}

func formatAcquirerTime(hostDate, hostTime string) (*time.Time, error) {
	if len(hostDate) != 4 {
		return nil, errors.New("host date format error")
	}
	if len(hostTime) != 6 {
		return nil, errors.New("host time format error")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%04d", time.Now().Year()))
	sb.WriteString(hostDate)
	sb.WriteString(hostTime)

	timeZone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("LoadLocation fail->", err.Error())
		return nil, err
	}
	ret, err := time.ParseInLocation("20060102150405", sb.String(), timeZone)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
