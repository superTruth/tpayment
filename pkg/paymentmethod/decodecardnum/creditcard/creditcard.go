package creditcard

import (
	"errors"
	"strconv"
	"strings"
	"tpayment/conf"
)

var Rules = []*brandRule{
	{
		CardBrand:  conf.Visa,
		CardNumLen: []int{13, 16, 19},
		CardNumPreFix: []*preFixNum{
			{
				Start: 400000,
				End:   500000,
			},
		},
	},
	{
		CardBrand:  conf.MasterCard,
		CardNumLen: []int{16},
		CardNumPreFix: []*preFixNum{
			{
				Start: 222100,
				End:   273000,
			},
			{
				Start: 510000,
				End:   560000,
			},
		},
	},
	{
		CardBrand:  conf.UnionPay,
		CardNumLen: []int{16, 17, 18, 19},
		CardNumPreFix: []*preFixNum{
			{
				Start: 620000,
				End:   630000,
			},
		},
	},
	{
		CardBrand:  conf.AE,
		CardNumLen: []int{15},
		CardNumPreFix: []*preFixNum{
			{
				Start: 340000,
				End:   350000,
			},
			{
				Start: 370000,
				End:   380000,
			},
		},
	},
	{
		CardBrand:  conf.JCB,
		CardNumLen: []int{16},
		CardNumPreFix: []*preFixNum{
			{
				Start: 352800,
				End:   359000,
			},
		},
	},
}

type brandRule struct {
	CardBrand     string
	CardNumLen    []int
	CardNumPreFix []*preFixNum
}

type preFixNum struct {
	Start int
	End   int
}

func Decode(cardNum string) (string, error) {
	cardLen := len(strings.TrimSpace(cardNum))
	numPreFix6, _ := strconv.Atoi(cardNum[:6])

	for _, rule := range Rules {
		// 匹配卡长度
		cardLenCompare := false
		for _, cardLenLoop := range rule.CardNumLen {
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
		for _, cardNumPreFixLoop := range rule.CardNumPreFix {
			if numPreFix6 >= cardNumPreFixLoop.Start &&
				numPreFix6 < cardNumPreFixLoop.End {
				preFixNumMatch = true
				break
			}
		}
		if !preFixNumMatch {
			continue
		}
		return rule.CardBrand, nil
	}

	return "", errors.New("not support card brand")
}
