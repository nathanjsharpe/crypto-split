package main

import (
	"flag"
	"fmt"
	"github.com/nathanjsharpe/crypto-split/coinbase"
	"github.com/nathanjsharpe/crypto-split/internal/app"
	"os"
)

type config struct {
	fiat        string
	holdingsAmt float64
	crypto1     string
	crypto2     string
}

func usage() {
	fmt.Println("Usage: crypto-split <amount> <crypto currency 1> <crypto currency 2>")
	fmt.Println("You will be given how much of the given crypto currencies to buy to have a split of 70% crypto currency 1 and 30% crypto currency 2.")
	fmt.Println("For example, if you have $100 and want a 70/30 split in BTC and ETH, you would run: crypto-split 100 BTC ETH")
	flag.PrintDefaults()
}

func main() {
	var cfg app.Config

	flag.StringVar(&cfg.Fiat, "fiat", "USD", "The fiat currency to use for holdings, as ISO 4217 currency code")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	cfg.Client = coinbase.NewClient()
	cfg.Out = os.Stdout

	err := app.Run(&cfg, flag.Args())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
