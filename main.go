package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nathanjsharpe/crypto-split/app"
	"github.com/nathanjsharpe/crypto-split/coinbase"
	"os"
	"strings"
)

var fiat string

func init() {
	flag.StringVar(&fiat, "fiat", "USD", "The fiat currency to use for holdings, as ISO 4217 currency code")
}

func usage() {
	fmt.Println("Usage: crypto-split <amount> <crypto currency 1> <crypto currency 2>")
	fmt.Println("You will be given how much of the given crypto currencies to buy to have a split of 70% crypto currency 1 and 30% crypto currency 2.")
	fmt.Println("For example, if you have $100 and want a 70/30 split in BTC and ETH, you would run: crypto-split 100 BTC ETH")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	switch flag.NArg() {
	case 0:
		flag.Usage()
		os.Exit(0)
	case 3:
	default:
		fmt.Println("Please provide 3 arguments.")
		flag.Usage()
		os.Exit(1)
	}

	a := app.Application{
		Fiat:   fiat,
		Client: coinbase.NewClient(),
	}

	s, err := splits(&a, flag.Args())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(strings.Join(s, ""))
}

func splits(a *app.Application, args []string) ([]string, error) {
	err := a.ParsePosArgs(args)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid arguments: %v", err))
	}

	return a.BuyInstructions()
}
