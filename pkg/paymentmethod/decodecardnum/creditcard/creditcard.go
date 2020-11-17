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
				end:   499999,
			},
		},
	},
	{
		cardBrand:  conf.MasterCard,
		cardNumLen: []int{16},
		cardNumPreFix: []*preFixNum{
			{
				start: 222100,
				end:   272999,
			},
			{
				start: 510000,
				end:   559999,
			},
		},
	},
	{
		cardBrand:  conf.UnionPay,
		cardNumLen: []int{16, 17, 18, 19},
		cardNumPreFix: []*preFixNum{
			{
				start: 620000,
				end:   629999,
			},
		},
	},
	{
		cardBrand:  conf.AE,
		cardNumLen: []int{15},
		cardNumPreFix: []*preFixNum{
			{
				start: 340000,
				end:   349999,
			},
			{
				start: 370000,
				end:   379999,
			},
		},
	},
	{
		cardBrand:  conf.JCB,
		cardNumLen: []int{16},
		cardNumPreFix: []*preFixNum{
			{
				start: 352800,
				end:   358999,
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
	//fmt.Println("cardNum->", cardNum, "card Len->", cardLen, "card prefix->", numPreFix6)

	for _, rule := range rules {
		//fmt.Println("compare ->", rule.cardBrand, "===================")
		// 匹配卡长度
		cardLenCompare := false
		for _, cardLenLoop := range rule.cardNumLen {
			//fmt.Println("match len ", cardLenLoop, "->", cardLen == cardLenLoop)
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
			//fmt.Println("match prefix start:", cardNumPreFixLoop.start, " end:", cardNumPreFixLoop.end)
			if numPreFix6 >= cardNumPreFixLoop.start &&
				numPreFix6 <= cardNumPreFixLoop.end {
				//fmt.Println("true")
				preFixNumMatch = true
				break
			}
			//fmt.Println("false")
		}
		if !preFixNumMatch {
			continue
		}
		return rule.cardBrand, nil
	}

	return "", errors.New("not support card brand")
}
