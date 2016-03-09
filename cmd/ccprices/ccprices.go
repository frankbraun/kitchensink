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

const (
	euroAPI  = "http://api.fixer.io/latest"
	xauAPI   = "https://www.quandl.com/api/v3/datasets/LBMA/GOLD.json?limit=1"
	xagAPI   = "https://www.quandl.com/api/v3/datasets/LBMA/SILVER.json?limit=1"
	coinsAPI = "http://coinmarketcap.northpole.ro/api/v5/all.json"
)

var (
	// Quandl API key can be set via environment variable QUANDL_API_KEY
	quandl = os.Getenv("QUANDL_API_KEY")
	coins  = []string{"Bitcoin", "Decred", "Namecoin"}
)

type result struct {
	symbol string
	price  float64
}

func getEuroExchangeRates() (map[string]interface{}, error) {
	resp, err := http.Get(euroAPI)
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	jsn := make(map[string]interface{})
	if err := json.Unmarshal(b, &jsn); err != nil {
		return nil, err
	}
	return jsn["rates"].(map[string]interface{}), nil
}

func getLBMAPrice(api string, dataIndex int) (float64, error) {
	if quandl != "" {
		api += "?api_key=" + quandl
	}
	resp, err := http.Get(api)
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0, err
	}
	jsn := make(map[string]interface{})
	if err := json.Unmarshal(b, &jsn); err != nil {
		return 0, err
	}
	data := jsn["dataset"].(map[string]interface{})["data"].([]interface{})
	var price float64
	if data[0].([]interface{})[dataIndex] != nil {
		// p.m. price is available
		price = data[0].([]interface{})[dataIndex].(float64)
	} else {
		// p.m. price is not available, use a.m. price instead
		price = data[0].([]interface{})[dataIndex-1].(float64)
	}
	return price, nil
}

func getCoinPrices() ([]interface{}, error) {
	resp, err := http.Get(coinsAPI)
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
	// get euro exchange rates
	rates, err := getEuroExchangeRates()
	if err != nil {
		fatal(err)
	}
	// get gold price
	xau, err := getLBMAPrice(xauAPI, 6)
	if err != nil {
		fatal(err)
	}
	// get silver price
	xag, err := getLBMAPrice(xagAPI, 3)
	if err != nil {
		fatal(err)
	}
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
	// output all prices
	t := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("P %s USD %11.6f EUR\n", t, 1/rates["USD"].(float64))
	fmt.Printf("P %s GBP %11.6f EUR\n", t, 1/rates["GBP"].(float64))
	fmt.Printf("P %s CHF %11.6f EUR\n", t, 1/rates["CHF"].(float64))
	fmt.Printf("P %s CZK %11.6f EUR\n", t, 1/rates["CZK"].(float64))
	fmt.Printf("P %s XAU %11.6f EUR\n", t, xau)
	fmt.Printf("P %s XAG %11.6f EUR\n", t, xag)
	for _, name := range coins {
		fmt.Printf("P %s %s %11.6f EUR\n", t, prices[name].symbol,
			prices[name].price)
	}
}
