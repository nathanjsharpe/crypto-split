package app

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Fiat   string
	Client CryptoClient
	Out    io.Writer

	holdings float64
	crypto1  string
	crypto2  string
}

type CryptoClient interface {
	ExchangeRate(fiat string, crypto string) (float64, error)
}

func Run(cfg *Config, args []string) error {
	if len(args) != 3 {
		return errors.New("you must provide 3 arguments")
	}

	amtInt, err := strconv.Atoi(args[0])
	if err != nil || amtInt <= 0 {
		return errors.New("first argument must be a positive integer")
	}

	cfg.holdings = float64(amtInt)
	cfg.crypto1 = strings.ToUpper(args[1])
	cfg.crypto2 = strings.ToUpper(args[2])

	p, err := purchases(cfg, map[string]float64{cfg.crypto1: 0.7, cfg.crypto2: 0.3})
	if err != nil {
		return err
	}

	for _, d := range p {
		_, err = fmt.Fprintf(cfg.Out, "%.2f %v => %.5f %v\n", d.fiatAmt, cfg.Fiat, d.cryptoAmt, d.curr)
		if err != nil {
			return err
		}
	}

	return nil
}

type cryptoPurchase struct {
	fiatAmt   float64
	cryptoAmt float64
	curr      string
}

func purchases(cfg *Config, distribution map[string]float64) ([]cryptoPurchase, error) {
	var results []cryptoPurchase

	for curr, amt := range distribution {
		r, err := cfg.Client.ExchangeRate(cfg.Fiat, curr)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to find exchange rate for %v => %v", cfg.Fiat, curr))
		}
		results = append(results, cryptoPurchase{
			fiatAmt:   cfg.holdings * amt,
			cryptoAmt: cfg.holdings * amt * r,
			curr:      curr,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].fiatAmt > results[j].fiatAmt
	})

	return results, nil
}
