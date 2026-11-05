package payments

import (
	"testing"

	"paypilot/internal/paypal"
)

func TestProviderCaptureMatchesExpectedAmountAndCurrency(t *testing.T) {
	tests := []struct {
		name   string
		amount paypal.Money
		want   bool
	}{
		{name: "matching amount and currency", amount: paypal.Money{Value: "25.00", CurrencyCode: "USD"}, want: true},
		{name: "currency comparison is case insensitive", amount: paypal.Money{Value: "25.00", CurrencyCode: "usd"}, want: true},
		{name: "amount differs", amount: paypal.Money{Value: "24.99", CurrencyCode: "USD"}},
		{name: "currency differs", amount: paypal.Money{Value: "25.00", CurrencyCode: "EUR"}},
		{name: "missing amount", amount: paypal.Money{CurrencyCode: "USD"}},
		{name: "malformed amount", amount: paypal.Money{Value: "25.001", CurrencyCode: "USD"}},
		{name: "zero amount", amount: paypal.Money{Value: "0.00", CurrencyCode: "USD"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := providerCaptureMatches(test.amount, 2500, "USD"); got != test.want {
				t.Fatalf("providerCaptureMatches(%+v) = %t, want %t", test.amount, got, test.want)
			}
		})
	}
}
