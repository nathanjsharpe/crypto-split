package app

import (
	"errors"
	"fmt"
	"strconv"
)

type Application struct {
	Fiat       string
	Holdings   float64
	CryptoDist map[string]float64
	Client     CryptoClient
}

type CryptoClient interface {
	ExchangeRate(fiat string, crypto string) (float64, error)
}

func (a *Application) ParsePosArgs(args []string) error {
	amtInt, err := strconv.Atoi(args[0])
	if err != nil || amtInt <= 0 {
		return errors.New("invalid amount - must be a positive integer")
	}

	a.Holdings = float64(amtInt)
	a.CryptoDist = map[string]float64{
		args[1]: 0.7,
		args[2]: 0.3,
	}

	return nil
}

func (a *Application) BuyInstructions() ([]string, error) {
	var results []string

	for curr, amt := range a.CryptoDist {
		r, err := a.Client.ExchangeRate(a.Fiat, curr)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to find exchange rate for %v => %v", a.Fiat, curr))
		}
		results = append(results, fmt.Sprintf("$%.2f => %.4f %v", a.Holdings*amt, a.Holdings*amt*r, curr))
	}

	return results, nil
}
