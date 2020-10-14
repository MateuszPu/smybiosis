package payprovider

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrencyTwoDecimal(t *testing.T) {
	assertion := assert.New(t)
	pln := AllCurrencies()["pln"]
	payment, fee := pln.CalculatePayment(15.56, 0.005)

	assertion.Equal(int64(1556), payment)
	assertion.Equal(int64(7), fee)
}

func TestCurrencyZeroDecimal(t *testing.T) {
	assertion := assert.New(t)
	parameters := []string{
		"bif",
		"clp",
		"djf",
		"gnf",
		"jpy",
		"kmf",
		"krw",
		"mga",
		"pyg",
		"rwf",
		"ugx",
		"vnd",
		"vuv",
		"xaf",
		"xof",
		"xpf",
	}
	for _, param := range parameters {
		t.Run("Verify address and name", func(t *testing.T) {
			currency := AllCurrencies()[param]
			payment, fee := currency.CalculatePayment(1509.55, 0.005)
			assertion.Equal(int64(1509), payment, fmt.Sprintf("Payment should be equals 15 for non decimall currency %s", param))
			assertion.Equal(int64(7), fee, fmt.Sprintf("Fee should be equals 7 for non decimall currency %s", param))
		})
	}
}
