package app

import (
	"errors"
	"fmt"
	"strconv"
)

type Application struct {
	Fiat     string
	Holdings float64
	Crypto1  string
	Crypto2  string
	Client   CryptoClient
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
	a.Crypto1 = args[1]
	a.Crypto2 = args[2]

	return nil
}

func (a *Application) PrintSplit() error {
	rate1, err := a.Client.ExchangeRate(a.Fiat, a.Crypto1)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to find exchange rate for %v => %v", a.Fiat, a.Crypto1))
	}
	rate2, err := a.Client.ExchangeRate(a.Fiat, a.Crypto2)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to find exchange rate for %v => %v", a.Fiat, a.Crypto2))
	}

	fmt.Printf("$%.2f => %.5f %v\n", a.Holdings*0.7, a.Holdings*0.7*rate1, a.Crypto1)
	fmt.Printf("$%.2f => %.5f %v\n", a.Holdings*0.3, a.Holdings*0.3*rate2, a.Crypto2)

	return nil
}
