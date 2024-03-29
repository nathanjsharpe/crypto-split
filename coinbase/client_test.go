package coinbase

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_ExchangeRate(t *testing.T) {
	type args struct {
		fiat   string
		crypto string
	}
	type call struct {
		args args
		want float64
	}
	type err struct {
		want bool
		msg  string
	}
	tests := []struct {
		name  string
		calls []call
		reqs  int
		err   err
	}{
		{
			"USD to BTC, USD to BTC",
			[]call{
				{args{"USD", "BTC"}, 0.0123},
				{args{"USD", "ETH"}, 0.234},
			},
			1,
			err{false, ""},
		},
		{
			"EUR to BTC, EUR to BTC",
			[]call{
				{args{"EUR", "BTC"}, 0.0321},
				{args{"EUR", "ETH"}, 0.432},
			},
			1,
			err{false, ""},
		},
		{
			"USD to BTC, EUR to BTC",
			[]call{
				{args{"USD", "BTC"}, 0.0123},
				{args{"EUR", "BTC"}, 0.0321},
			},
			2,
			err{false, ""},
		},
		{
			"USD to BTC, EUR to BTC, USD to ETH, EUR to ETH",
			[]call{
				{args{"USD", "BTC"}, 0.0123},
				{args{"EUR", "BTC"}, 0.0321},
				{args{"USD", "ETH"}, 0.234},
				{args{"EUR", "ETH"}, 0.432},
			},
			2,
			err{true, ""},
		},
		{
			"Unknown fiat currency",
			[]call{
				{args{"WAT", "BTC"}, 0},
			},
			1,
			err{true, "no exchange rate"},
		},
		{
			"Unknown crypto currency",
			[]call{
				{args{"USD", "WAT"}, 0},
			},
			1,
			err{true, "no exchange rate"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqs := 0
			// We set up a test server to control the responses and avoid unnecessary http calls to an external
			// server. The server responds with json in the same shape as the coinbase API.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqs++
				if reqs > tt.reqs {
					t.Errorf("expected %v requests, got extra request: %v", tt.reqs, r.URL)
				}
				w.Header().Set("Content-Type", "application/json")
				switch fiat := r.URL.Query()["currency"][0]; fiat {
				case "USD":
					fmt.Fprint(w, `{
						"data": {
							"currency": "USD",
							"rates": {
								"BTC": "0.0123",
								"ETH": "0.234"
							}
						}
					}`)
				case "EUR":
					fmt.Fprint(w, `{
						"data": {
							"currency": "EUR",
							"rates": {
								"BTC": "0.0321",
								"ETH": "0.432"
							}
						}
					}`)
				default:
					fmt.Fprintf(w, `{
						"data": {
							"currency": "%v",
							"rates": {
								"%v": "1.0"
							}
						}
					}`, fiat, fiat)
				}
			}))
			defer ts.Close()

			c := NewClient()
			c.baseUrl = ts.URL
			for _, call := range tt.calls {
				got, err := c.ExchangeRate(call.args.fiat, call.args.crypto)
				if err != nil {
					if !tt.err.want {
						t.Errorf("expected no error, got %v", err)
					}
					if !strings.Contains(err.Error(), tt.err.msg) {
						t.Errorf("expected error message to contain '%v', got '%v'", tt.err.msg, err)
					}
				}
				if got != call.want {
					t.Errorf("expected %v, got %v", call.want, got)
				}
			}
		})
	}
}

// This test will only run if the environment variable CRYPTO_SPLIT_SEND_EXTERNAL_REQUESTS is set to "true".
// It makes an actual request to the coinbase API, so it is not run by default. It can act as a sort of once per release
// smoke test to make sure the API has not changed.
func Test_ExchangeRate_real(t *testing.T) {
	v, ok := os.LookupEnv("CRYPTO_SPLIT_SEND_EXTERNAL_REQUESTS")
	if !ok || v != "true" {
		t.Skip("skipping test; set CRYPTO_SPLIT_SEND_EXTERNAL_REQUESTS to run")
	}

	c := NewClient()
	r, err := c.ExchangeRate("USD", "BTC")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if r == 0 {
		t.Errorf("expected rate to be non-zero, got %v", r)
	}
}
