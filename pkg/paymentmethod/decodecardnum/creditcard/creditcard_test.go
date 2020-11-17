package creditcard

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	testCard := []string{
		"5599110224911916",
		"5599110224922916",
		"4514617642433897",
		"5599110224944916",
		"5599110224955916",
		"5599110224966916",
		"5599110224977916",
		"5599110224988916",
		"4514617642499897",
		"5599110224911916",
		"5599110224922916",
		"5599110224933916",
		"4514617642433897",
		"5599110224933916",
		"5599110224933916",
		"4514617642433897",
		"4514617642433897",
		"4514617642433897",
		"4514617642433897",
		"3528110224933916",
		"5599110224933916",
		"5599110224933916",
		"349911022493391",
		"5599110224933916",
		"6214617642433897"}

	for _, cardNum := range testCard {
		brand, err := Decode(cardNum)

		if err != nil {
			t.Error(err.Error())
			return
		}

		fmt.Println("card number: ", cardNum, " ->", brand)
	}

}
