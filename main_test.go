package main

import (
	"errors"
	"fmt"
	"github.com/nathanjsharpe/crypto-split/app"
	"strings"
	"testing"
)

type testClient struct{}

// We'll use a test client that implements the CryptoClient interface. It supports just enough currencies to cover test
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

func Test_splits_errors(t *testing.T) {
	type args struct {
		a    *app.Application
		args []string
	}
	tests := []struct {
		name string
		args args
		err  string
	}{
		{
			"With a non-number amount",
			args{nil, []string{"abc", "BTC", "ETH"}},
			"positive integer",
		},
		{
			"With a float amount",
			args{nil, []string{"1.5", "BTC", "ETH"}},
			"positive integer",
		},
		{
			"With a negative integer amount",
			args{nil, []string{"-100", "BTC", "ETH"}},
			"positive integer",
		},
		{
			"with zero amount",
			args{nil, []string{"0", "BTC", "ETH"}},
			"positive integer",
		},
		{
			"With an unsupported fiat currency",
			args{&app.Application{Fiat: "WAT", Client: &testClient{}}, []string{"100", "BTC", "ETH"}},
			"failed to find exchange rate",
		},
		{
			"With an unsupported crypto currency",
			args{&app.Application{Fiat: "USD", Client: &testClient{}}, []string{"100", "WAT", "ETH"}},
			"failed to find exchange rate",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := splits(tt.args.a, tt.args.args)
			if err == nil {
				t.Errorf("expected error matcing '%v', but no error", tt.err)
			}
			if !strings.Contains(fmt.Sprint(err), tt.err) {
				t.Errorf("expected error matching '%v', but got '%v'", tt.err, err)
			}
		})
	}
}

func Test_splits_success(t *testing.T) {
	a := &app.Application{Fiat: "USD", Client: &testClient{}}
	tests := []struct {
		name    string
		args    []string
		results []string
	}{
		{
			"With 100 as the amount",
			[]string{"100", "BTC", "ETH"},
			[]string{"70.00 USD => 70.0000 BTC", "30.00 USD => 60.0000 ETH"},
		},
		{
			"With 2 as the amount",
			[]string{"1", "BTC", "ETH"},
			[]string{"0.70 USD => 0.7000 BTC", "0.30 USD => 0.6000 ETH"},
		},
		{
			"With 7 as the amount",
			[]string{"7", "BTC", "ETH"},
			[]string{"4.90 USD => 4.9000 BTC", "2.10 USD => 4.2000 ETH"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := splits(a, tt.args)
			if err != nil {
				t.Errorf("expected no error, but received %v", err)
			}
			for _, r := range s {
				match := false
				for _, expected := range tt.results {
					if r == expected {
						match = true
						break
					}
				}
				if !match {
					t.Errorf("expected result %v, but got %v", tt.results, s)
				}
			}

		})
	}
}
