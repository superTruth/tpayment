package qrcode

import (
	"errors"
	"strconv"
	"strings"
	"tpayment/conf"
)

var rules = []*brandRule{
	{
		cardBrand:  conf.WeChatPay,
		cardNumLen: []int{17},
		cardNumPreFix: []*preFixNum{
			{
				start: 100000,
				end:   200000,
			},
		},
	},
	{
		cardBrand:  conf.Alipay,
		cardNumLen: []int{16},
		cardNumPreFix: []*preFixNum{
			{
				start: 280000,
				end:   290000,
			},
		},
	},
	{
		cardBrand:  conf.UnionPay,
		cardNumLen: []int{16, 17, 18, 19},
		cardNumPreFix: []*preFixNum{
			{
				start: 620000,
				end:   630000,
			},
		},
	},
}

type brandRule struct {
	cardBrand     string
	cardNumLen    []int
	cardNumPreFix []*preFixNum
}

type preFixNum struct {
	start int
	end   int
}

func Decode(cardNum string) (string, error) {
	cardLen := len(strings.TrimSpace(cardNum))
	numPreFix6, _ := strconv.Atoi(cardNum[:6])

	for _, rule := range rules {
		// 匹配卡长度
		cardLenCompare := false
		for _, cardLenLoop := range rule.cardNumLen {
			if cardLen == cardLenLoop {
				cardLenCompare = true
				break
			}
		}

		if !cardLenCompare {
			continue
		}

		// 匹配前面6个数据
		preFixNumMatch := false
		for _, cardNumPreFixLoop := range rule.cardNumPreFix {
			if numPreFix6 >= cardNumPreFixLoop.start &&
				numPreFix6 < cardNumPreFixLoop.end {
				preFixNumMatch = true
				break
			}
		}
		if !preFixNumMatch {
			continue
		}
		return rule.cardBrand, nil
	}

	return "", errors.New("not support card brand")
}
