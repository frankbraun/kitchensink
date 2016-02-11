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

var coins = []string{"BTC", "DCR", "NMC"}

type result struct {
	symbol string
	price  float64
	err    error
}

func getCoinPrice(ch chan<- *result, sym string) {
	resp, err := http.Get("http://coinmarketcap.northpole.ro/api/v5/" + sym + ".json")
	if err != nil {
		ch <- &result{err: err}
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- &result{err: err}
		return
	}
	jsn := make(map[string]interface{})
	if err := json.Unmarshal(b, &jsn); err != nil {
		ch <- &result{err: err}
		return
	}
	s, _ := json.MarshalIndent(jsn, "", "  ")
	fmt.Println(string(s))
	f := jsn["price"].(map[string]interface{})["eur"].(string)
	p, err := strconv.ParseFloat(f, 64)
	if err != nil {
		ch <- &result{err: err}
		return
	}
	ch <- &result{symbol: sym, price: p}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	ch := make(chan *result)
	for _, sym := range coins {
		go getCoinPrice(ch, sym)
	}
	prices := make(map[string]float64)
	for range coins {
		r := <-ch
		if r.err != nil {
			fatal(r.err)
		}
		prices[r.symbol] = r.price
	}
	t := time.Now().Format("2006/01/02 15:04:05")
	for _, sym := range coins {
		fmt.Printf("P %s %s %7.2f EUR\n", t, sym, prices[sym])
	}
}
