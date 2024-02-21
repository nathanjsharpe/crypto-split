package coinbase

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type exchangeRatesResp struct {
	Data struct {
		Rates map[string]string `json:"rates"`
	} `json:"data"`
}

type RateCache struct {
	rates     map[string]string
	fetchedAt time.Time
}

func (cache *RateCache) rate(curr string) (float64, error) {
	rateStr, ok := cache.rates[curr]
	if !ok {
		return 0, errors.New("no exchange rate")
	}

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return 0, errors.New("failed to parse exchange rate")
	}

	return rate, nil
}

func fetchRates(fiat string, c *Client) (map[string]string, error) {
	body, err := c.get("/exchange-rates", map[string]string{"currency": fiat})
	var cryptoResp exchangeRatesResp
	err = json.Unmarshal(body, &cryptoResp)
	if err != nil {
		return nil, errors.New("failed to parse json")
	}

	return cryptoResp.Data.Rates, nil
}
