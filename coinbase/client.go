package coinbase

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	cachedRates map[string]RateCache
}

const baseUrl string = "https://api.coinbase.com/v2"

func NewClient() *Client {
	return &Client{
		cachedRates: make(map[string]RateCache),
	}
}

func (c *Client) ExchangeRate(fiat string, crypto string) (float64, error) {
	_, cached := c.cachedRates[fiat]
	if !cached {
		rates, err := fetchRates(fiat)
		if err != nil {
			return 0, err
		}
		c.cachedRates[fiat] = RateCache{
			rates:     rates,
			fetchedAt: time.Now(),
		}
	}

	cache := c.cachedRates[fiat]

	return cache.rate(crypto)
}

func get(path string, params map[string]string) ([]byte, error) {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	for _, param := range params {
		q.Add(param, params[param])
	}
	base.RawQuery = q.Encode()

	base.Path += path
	resp, err := http.Get(base.String())
	if err != nil {
		return nil, errors.New("error fetching exchange rates")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to parse http response body")
	}

	return body, nil
}
