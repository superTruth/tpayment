package format_utils

import (
	"errors"
	"strings"
)

func FloatAmount2Int(src string, decimalLen int) (string, error) {
	src = DeleteAppendString(src, true, '0') // 去除掉开头的0
	amountSp := strings.Split(src, ".")
	if len(amountSp) > 2 {
		return "", errors.New("source amount error->" + src)
	}
	decimalPart := ""
	if len(amountSp) == 2 {
		decimalPart = amountSp[1]
	}

	decimalPart = AppendString(decimalPart, decimalLen, false, '0')

	ret := amountSp[0] + decimalPart
	ret = DeleteAppendString(ret, true, '0') // 去除掉开头的0
	if len(ret) == 0 {
		ret = "0"
	}

	return ret, nil
}
