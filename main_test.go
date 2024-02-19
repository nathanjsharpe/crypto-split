package main

import (
	"errors"
	"fmt"
	"github.com/nathanjsharpe/crypto-split/app"
	"regexp"
	"testing"
)

type testClient struct{}

func (c *testClient) ExchangeRate(fiat, crypto string) (float64, error) {
	switch {
	case fiat == "USD" && crypto == "BTC":
		return 1, nil
	case fiat == "USD" && crypto == "ETH":
		return 2, nil
	case fiat == "USD":
		return 0, errors.New("unsupported crypto currency")
	default:
		return 0, errors.New("unsupported fiat currency")
	}
}

func Test_execute_errors(t *testing.T) {
	type args struct {
		a    *app.Application
		args []string
	}
	type err struct {
		wantErr bool
		msg     string
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
			err := execute(tt.args.a, tt.args.args)
			if err == nil {
				t.Errorf("expected error matcing '%v', but no error", tt.err)
			}
			if match, _ := regexp.MatchString(tt.err, fmt.Sprint(err)); !match {
				t.Errorf("expected error matching '%v', but got '%v'", tt.err, err)
			}
		})
	}
}
