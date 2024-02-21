package app

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type Application struct {
	Fiat       string
	Holdings   float64
	CryptoDist map[string]float64
	Client     CryptoClient
}

type CryptoPurchase struct {
	fiatAmt   float64
	cryptoAmt float64
	curr      string
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
	dirs, err := a.splitBuys()
	if err != nil {
		return nil, err
	}

	var results []string

	for _, d := range dirs {
		results = append(results, fmt.Sprintf("%.2f %v => %.4f %v", d.fiatAmt, a.Fiat, d.cryptoAmt, d.curr))
	}

	return results, nil
}

func (a *Application) splitBuys() ([]CryptoPurchase, error) {
	var results []CryptoPurchase

	for curr, amt := range a.CryptoDist {
		r, err := a.Client.ExchangeRate(a.Fiat, curr)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to find exchange rate for %v => %v", a.Fiat, curr))
		}
		results = append(results, CryptoPurchase{
			fiatAmt:   a.Holdings * amt,
			cryptoAmt: a.Holdings * amt * r,
			curr:      curr,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].fiatAmt > results[j].fiatAmt
	})

	return results, nil
}
