This is a small utility to help intrepid investors buy cryptocurrencies in a 70/30 split. 

Given an amount of money to invest and a couple of cryptocurrencies to buy, it uses the Coinbase API to find exchange rates and tells the user how much of each cryptocurrency to buy to invest 70% of the investment in the first cryptocurrency and 30% in the second.

It supports any fiat currency and cryptocurrency that Coinbase does. You can get an idea of basic usage by just running without arguments.

## Overview

There are three packages:

- `coinbase` handles communication with the Coinbase API: fetching exchange rates, parsing JSON, etc.
- `app` handles the core logic of the app: figuring out which currencies to use and calculating the amount of cryptocurrencies to buy.
- `main` ties these other two modules together in a way that prevents them from needing to know details about one another and displays results and errors to the user.

We keep `app` and `coinbase` separate by defining an interface for a crypto client in `app` then injecting the Coinbase client when using `app` from `main`. This allows working on `app` and `coinbase` in relative isolation and would allow easily using a different API should Coinbase disappear (not that companies involved in crypto ever disappear!).

The Coinbase client has its own `baseUrl`, but that can be overwritten (within the package) for easier testing.

This setup should make it easier to add more features, like defining custom splits. I considered doing that here, but I had already added support for additional fiat currencies, and a line must be drawn somewhere.

There are probably more idiomatic ways to do what I tried to do, but I did my best to seek out conventions and use them.

## Testing

The functionality of the `app` and `main` packages are both covered by specs for `main`. There are separate specs for the `coinbase` package that test the `client`. These tests cover the functionality in `rates`.

The `main` specs define a crypto client that implements the CryptoClient interface and uses that for testing. The `coinbase` tests similarly define a test http server and make requests against that. There is one test that does hit the actual Coinbase API, but it requires setting an environment variable to run. That can give us some basic protection against the API changing or just disappearing.

## Prompt

Your Task:
You are to make a cli that takes in a USD amount as holdings, and calculates the 70/30 split for 2 given cryptocurrencies. Stated simply: I have $X I want to keep in BTC and ETH, 70/30 split. How many of each should I buy? An example usage would look like:

binary_name 100 BTC ETH

This output should be in the following format:

$70.00 => 0.0025 BTC
$30.00 => 0.0160 ETH

This output tells us: Of our 100$ holdings, 70% of that is 70$, which buys 0.0025 BTC, and 30% of our holdings is 30$, which buys 0.016 ETH