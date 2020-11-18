package creditcard

import (
	"errors"
	"strconv"
	"strings"
	"tpayment/conf"
)

var rules = []*brandRule{
	{
		cardBrand:  conf.Visa,
		cardNumLen: []int{13, 16, 19},
		cardNumPreFix: []*preFixNum{
			{
				start: 400000,
				end:   500000,
			},
		},
	},
	{
		cardBrand:  conf.MasterCard,
		cardNumLen: []int{16},
		cardNumPreFix: []*preFixNum{
			{
				start: 222100,
				end:   273000,
			},
			{
				start: 510000,
				end:   560000,
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
	{
		cardBrand:  conf.AE,
		cardNumLen: []int{15},
		cardNumPreFix: []*preFixNum{
			{
				start: 340000,
				end:   350000,
			},
			{
				start: 370000,
				end:   380000,
			},
		},
	},
	{
		cardBrand:  conf.JCB,
		cardNumLen: []int{16},
		cardNumPreFix: []*preFixNum{
			{
				start: 352800,
				end:   359000,
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