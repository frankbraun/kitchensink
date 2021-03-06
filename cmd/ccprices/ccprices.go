// ccprices prints current currency prices in ledger format.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	euroAPI  = "http://data.fixer.io/api/latest"
	xauAPI   = "https://www.quandl.com/api/v3/datasets/LBMA/GOLD.json?limit=1"
	xagAPI   = "https://www.quandl.com/api/v3/datasets/LBMA/SILVER.json?limit=1"
	coinsAPI = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
)

var (
	// CoinMarketCap API key can be set via environment variable COINMARKETCAP_API_KEY
	coinmarketcap = os.Getenv("COINMARKETCAP_API_KEY")
	// Fixer API key can be set via environment variable FIXER_API_KEY
	fixer = os.Getenv("FIXER_API_KEY")
	// Quandl API key can be set via environment variable QUANDL_API_KEY
	quandl = os.Getenv("QUANDL_API_KEY")
	coins  = []string{
		"Bitcoin",
		"Bitcoin Cash",
		"Bitcoin Gold",
		"Bitcoin SV",
		"Dash",
		"Decred",
		"Ethereum",
		"Grin",
		"Litecoin",
		"Monero",
		"Particl",
		"Tezos",
		"Zcash",
	}
)

type result struct {
	symbol string
	price  float64
}

func httpGetWithWarning(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		warning(fmt.Sprintf("GET %s: %s", url, resp.Status))
		return nil, nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

func getEuroExchangeRates(api string) (map[string]interface{}, error) {
	if fixer != "" {
		api += "?access_key=" + fixer
	}
	b, err := httpGetWithWarning(api)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	jsn := make(map[string]interface{})
	if err := json.Unmarshal(b, &jsn); err != nil {
		return nil, err
	}
	return jsn["rates"].(map[string]interface{}), nil
}

func getLBMAPrice(api string, dataIndex int) (float64, error) {
	if quandl != "" {
		api += "&api_key=" + quandl
	}
	b, err := httpGetWithWarning(api)
	if err != nil {
		return 0, err
	}
	if b == nil {
		return 0, nil
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", coinsAPI, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "2000")
	q.Add("convert", "EUR")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", coinmarketcap)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		warning(fmt.Sprintf("GET %s: %s", coinmarketcap, resp.Status))
		return nil, nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsn := make(map[string]interface{})
	if err := json.Unmarshal(b, &jsn); err != nil {
		return nil, err
	}
	return jsn["data"].([]interface{}), nil
}

func warning(warn string) {
	fmt.Fprintf(os.Stderr, "%s: warning: %s\n", os.Args[0], warn)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	// get euro exchange rates
	rates, err := getEuroExchangeRates(euroAPI)
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
	var (
		names  map[string]struct{}
		prices map[string]*result
	)
	if all != nil {
		names = make(map[string]struct{})
		for _, name := range coins {
			names[name] = struct{}{}
		}
		prices = make(map[string]*result)
		// iterate over all coin informations
		var btc, bch float64
		for _, info := range all {
			coin := info.(map[string]interface{})
			name := coin["name"].(string)
			_, ok := names[name]
			if ok {
				// we are interested in this coin -> store price and symbol
				quote := coin["quote"].(map[string]interface{})
				eur := quote["EUR"].(map[string]interface{})
				p := eur["price"].(float64)
				prices[name] = &result{symbol: coin["symbol"].(string), price: p}
				if coin["symbol"] == "BTC" {
					btc = p
				}
				if coin["symbol"] == "BCH" {
					bch = p
				}
			}
		}
		fmt.Fprintf(os.Stderr, "BCH/BTC ratio: %.2f%%\n", bch*100.0/btc)
	}
	// output all prices
	t := time.Now().Format("2006/01/02 15:04:05")
	if rates != nil {
		fmt.Printf("P %s USD %11.6f EUR\n", t, 1/rates["USD"].(float64))
		fmt.Printf("P %s GBP %11.6f EUR\n", t, 1/rates["GBP"].(float64))
		fmt.Printf("P %s CHF %11.6f EUR\n", t, 1/rates["CHF"].(float64))
		fmt.Printf("P %s CZK %11.6f EUR\n", t, 1/rates["CZK"].(float64))
		fmt.Printf("P %s THB %11.6f EUR\n", t, 1/rates["THB"].(float64))
	}
	if xau != 0 {
		fmt.Printf("P %s XAU %11.6f EUR\n", t, xau)
	}
	if xag != 0 {
		fmt.Printf("P %s XAG %11.6f EUR\n", t, xag)
	}
	if all != nil {
		for _, name := range coins {
			price, ok := prices[name]
			if ok {
				fmt.Printf("P %s %s %11.6f EUR\n", t, price.symbol, price.price)
			} else {
				fmt.Fprintf(os.Stderr, "price for \"%s\" does not exist\n", name)
			}
		}
	}
}
