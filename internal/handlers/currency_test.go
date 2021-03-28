package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCurPair(t *testing.T) {
	cases := [...]struct {
		firstCur string
		secCur   string
		expPair  string
	}{
		{"доллар", "рубли", "USD_RUB"},
		{"доллара", "рубли", "USD_RUB"},
		{"долларов", "евро", "USD_EUR"},
		{"евро", "рубли", "EUR_RUB"},
		{"евро", "доллары", "EUR_USD"},
		{"рублей", "доллары", "RUB_USD"},
		{"рубль", "евро", "RUB_EUR"},
		{"рубля", "евро", "RUB_EUR"},
		{"доллар", "доллары", ""},
		{"рублей", "рубли", ""},
		{"евро", "евро", ""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("firstCur=%s;secCur=%s", c.firstCur, c.secCur), func(t *testing.T) {
			assert.Equal(t, getCurPair(c.firstCur, c.secCur), c.expPair)
		})
	}
}
