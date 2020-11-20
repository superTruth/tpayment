package creditcard

import (
	"errors"
	"strings"
)

func GetCardNumFromTrack2(track2 string) (string, error) {
	track2Format := strings.ToUpper(strings.TrimSpace(track2))
	track2Format = strings.ReplaceAll(track2Format, "=", "D")

	tk2SplitData := strings.Split(track2Format, "D")
	if len(tk2SplitData) < 2 {
		return "", errors.New("can't find split D")
	}

	// 去除非数字
	ret := strings.TrimFunc(tk2SplitData[0], func(r rune) bool {
		return r < '0' || r > '9'
	})

	return ret, nil
}
