package app

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type testClient struct{}

// We'll use a test client that implements the cryptoClient interface. It supports just enough currencies to cover test
// cases.
func (c *testClient) ExchangeRate(fiat, crypto string) (float64, error) {
	switch {
	case fiat == "USD" && crypto == "BTC":
		return 1.0, nil
	case fiat == "USD" && crypto == "ETH":
		return 2.0, nil
	case fiat == "USD":
		return 0, errors.New("unsupported crypto currency")
	default:
		return 0, errors.New("unsupported fiat currency")
	}
}

func TestRun(t *testing.T) {
	var b bytes.Buffer
	c := &Config{Fiat: "USD", Client: &testClient{}, Out: &b}

	type args struct {
		cfg  *Config
		args []string
	}
	tests := []struct {
		name    string
		args    args
		results []string
		err     string
	}{
		{
			"With a non-number amount",
			args{c, []string{"abc", "BTC", "ETH"}},
			nil,
			"positive integer",
		},
		{
			"With a float amount",
			args{c, []string{"1.5", "BTC", "ETH"}},
			nil,
			"positive integer",
		},
		{
			"With a negative integer amount",
			args{c, []string{"-100", "BTC", "ETH"}},
			nil,
			"positive integer",
		},
		{
			"With zero amount",
			args{c, []string{"0", "BTC", "ETH"}},
			nil,
			"positive integer",
		},
		{
			"With no arguments",
			args{c, []string{}},
			nil,
			"you must provide 3 arguments",
		},
		{
			"With too few arguments",
			args{c, []string{"100", "BTC"}},
			nil,
			"you must provide 3 arguments",
		},
		{
			"With too many arguments",
			args{c, []string{"100", "BTC", "ETH", "LTC"}},
			nil,
			"you must provide 3 arguments",
		},
		{
			"With an unsupported fiat currency",
			args{&Config{Fiat: "WAT", Client: &testClient{}}, []string{"100", "BTC", "ETH"}},
			nil,
			"failed to find exchange rate",
		},
		{
			"With an unsupported crypto currency",
			args{c, []string{"100", "WAT", "ETH"}},
			nil,
			"failed to find exchange rate",
		},
		{
			"With 100 as the amount",
			args{c, []string{"100", "BTC", "ETH"}},
			[]string{"70.00 USD => 70.00000 BTC", "30.00 USD => 60.00000 ETH"},
			"",
		},
		{
			"With 2 as the amount",
			args{c, []string{"1", "BTC", "ETH"}},
			[]string{"0.70 USD => 0.70000 BTC", "0.30 USD => 0.60000 ETH"},
			"",
		},
		{
			"With 7 as the amount",
			args{c, []string{"7", "BTC", "ETH"}},
			[]string{"4.90 USD => 4.90000 BTC", "2.10 USD => 4.20000 ETH"},
			"",
		},
		{
			"With lowercase crypto currencies",
			args{c, []string{"100", "btc", "eth"}},
			[]string{"70.00 USD => 70.00000 BTC", "30.00 USD => 60.00000 ETH"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b.Reset()
			err := Run(tt.args.cfg, tt.args.args)
			if tt.err != "" {
				switch {
				case err == nil:
					t.Errorf("expected error matcing '%v', but no error", tt.err)
				case !strings.Contains(fmt.Sprint(err), tt.err):
					t.Errorf("expected error matching '%v', but got '%v'", tt.err, err)
				}
			}
			if tt.err == "" && err != nil {
				t.Errorf("expected no error, but received %v", err)
			}
			if tt.results == nil {
				return
			}
			results := strings.Split(b.String(), "\n")
			for i, expected := range tt.results {
				if expected != results[i] {
					t.Errorf("expected result %v, but got %v", tt.results, results)
					return
				}
			}

		})
	}
}
