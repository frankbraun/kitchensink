// Copyright (c) 2016 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// ccprices prints current currency prices in ledger format.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var coins = []string{"Bitcoin", "Decred", "Namecoin"}

type result struct {
	symbol string
	price  float64
}

func getCoinPrices() ([]interface{}, error) {
	resp, err := http.Get("http://coinmarketcap.northpole.ro/api/v5/all.json")
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	jsn := make(map[string]interface{})
	if err := json.Unmarshal(b, &jsn); err != nil {
		return nil, err
	}
	return jsn["markets"].([]interface{}), nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	// get all coin prices
	all, err := getCoinPrices()
	if err != nil {
		fatal(err)
	}
	// construct map of coin names we are interested in
	names := make(map[string]struct{})
	for _, name := range coins {
		names[name] = struct{}{}
	}
	prices := make(map[string]*result)
	// iterate over all coin informations
	for _, info := range all {
		coin := info.(map[string]interface{})
		name := coin["name"].(string)
		_, ok := names[name]
		if ok {
			// we are interested in this coin -> store price and symbol
			f := coin["price"].(map[string]interface{})["eur"].(string)
			p, err := strconv.ParseFloat(f, 64)
			if err != nil {
				fatal(err)
			}
			prices[name] = &result{symbol: coin["symbol"].(string), price: p}
		}
	}
	// output coin prices
	t := time.Now().Format("2006/01/02 15:04:05")
	for _, name := range coins {
		fmt.Printf("P %s %s %7.2f EUR\n", t, prices[name].symbol,
			prices[name].price)
	}
}
