// ccprices prints current currency prices in ledger format.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	euroAPI      = "http://data.fixer.io/api/latest"
	xauAPI       = "https://www.quandl.com/api/v3/datasets/LBMA/GOLD.json?limit=1"
	xagAPI       = "https://www.quandl.com/api/v3/datasets/LBMA/SILVER.json?limit=1"
	coinsAPI     = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
	auroraAPI    = "https://api.coingecko.com/api/v3/coins/aurora-near"
	dsdAPI       = "https://api.coingecko.com/api/v3/simple/token_price/ethereum?contract_addresses=0xbd2f0cd039e0bfcf88901c98c0bfac5ab27566e3&vs_currencies=EUR"
	dsdContract  = "0xbd2f0cd039e0bfcf88901c98c0bfac5ab27566e3"
	esdAPI       = "https://api.coingecko.com/api/v3/simple/token_price/ethereum?contract_addresses=0x36f3fd68e7325a35eb768f1aedaae9ea0689d723&vs_currencies=EUR"
	esdContract  = "0x36f3fd68e7325a35eb768f1aedaae9ea0689d723"
	fraxAPI      = "https://api.coingecko.com/api/v3/simple/token_price/ethereum?contract_addresses=0x853d955acef822db058eb8505911ed77f175b99e&vs_currencies=EUR"
	fraxContract = "0x853d955acef822db058eb8505911ed77f175b99e"
	fxsAPI       = "https://api.coingecko.com/api/v3/simple/token_price/ethereum?contract_addresses=0x3432b6a60d23ca0dfca7761b7ab56459d9c964d0&vs_currencies=EUR"
	fxsContract  = "0x3432b6a60d23ca0dfca7761b7ab56459d9c964d0"
)

var (
	// CoinMarketCap API key can be set via environment variable COINMARKETCAP_API_KEY
	coinmarketcap = os.Getenv("COINMARKETCAP_API_KEY")
	// Fixer API key can be set via environment variable FIXER_API_KEY
	fixer = os.Getenv("FIXER_API_KEY")
	// Quandl API key can be set via environment variable QUANDL_API_KEY
	quandl = os.Getenv("QUANDL_API_KEY")
	coins  = []string{
		"Avalanche",
		"Balancer",
		"Bitcoin",
		"Bitcoin Cash",
		"Bitcoin Gold",
		"Bitcoin SV",
		"Curve DAO Token",
		"Dai",
		"Dash",
		"Decred",
		"DeFi Pulse Index",
		"Ethereum",
		"Fantom",
		"FTX Token",
		"Grin",
		"JOE",
		"Litecoin",
		"Marinade Staked SOL",
		"NEAR Protocol",
		"Nexo",
		"Monero",
		"PancakeSwap",
		"Particl",
		"Phonon DAO",
		"Polygon",
		"Raydium",
		"Saber",
		"Solana",
		"Stacks",
		"Terra",
		"Tezos",
		"USD Coin",
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
	if jsn["success"].(bool) == false {
		jsn, err := json.Marshal(jsn["error"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(os.Stderr, string(jsn))
		return nil, nil
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

func getAuroraPrice(api string) (map[string]interface{}, error) {
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
	marketData := jsn["market_data"].(map[string]interface{})
	currentPrice := marketData["current_price"].(map[string]interface{})
	return currentPrice, nil
}

func getDSDPrice(api string) (map[string]interface{}, error) {
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
	return jsn[dsdContract].(map[string]interface{}), nil
}

func getESDPrice(api string) (map[string]interface{}, error) {
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
	return jsn[esdContract].(map[string]interface{}), nil
}

func getFRAXPrice(api string) (map[string]interface{}, error) {
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
	return jsn[fraxContract].(map[string]interface{}), nil
}

func getFXSPrice(api string) (map[string]interface{}, error) {
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
	return jsn[fxsContract].(map[string]interface{}), nil
}

func getCoinPrices() ([]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", coinsAPI, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
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

	/*
		s, err := json.MarshalIndent(jsn, "", "  ")
		if err != nil {
			return nil, err
		}
		fmt.Println(string(s))
	*/

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
	// get AURORA price
	aurora, err := getAuroraPrice(auroraAPI)
	if err != nil {
		fatal(err)
	}
	// get DSD price
	dsd, err := getDSDPrice(dsdAPI)
	if err != nil {
		fatal(err)
	}
	// get ESD price
	esd, err := getESDPrice(esdAPI)
	if err != nil {
		fatal(err)
	}
	// get FRAX price
	frax, err := getFRAXPrice(fraxAPI)
	if err != nil {
		fatal(err)
	}
	// get FXS price
	fxs, err := getFXSPrice(fxsAPI)
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
		fmt.Printf("P %s HUF %11.6f EUR\n", t, 1/rates["HUF"].(float64))
		fmt.Printf("P %s THB %11.6f EUR\n", t, 1/rates["THB"].(float64))
		fmt.Printf("P %s CRC %11.6f EUR\n", t, 1/rates["CRC"].(float64))
		fmt.Printf("P %s AED %11.6f EUR\n", t, 1/rates["AED"].(float64))
		fmt.Printf("P %s MYR %11.6f EUR\n", t, 1/rates["MYR"].(float64))
	}
	if xau != 0 {
		fmt.Printf("P %s XAU %11.6f EUR\n", t, xau)
	}
	if xag != 0 {
		fmt.Printf("P %s XAG %11.6f EUR\n", t, xag)
	}
	var btc float64
	var eth float64
	if all != nil {
		for _, name := range coins {
			price, ok := prices[name]
			if ok {
				fmt.Printf("P %s %s %11.6f EUR\n", t, price.symbol, price.price)
				if price.symbol == "BTC" {
					btc = price.price
				}
				if price.symbol == "ETH" {
					eth = price.price
				}
			} else {
				fmt.Fprintf(os.Stderr, "price for \"%s\" does not exist\n", name)
			}
		}
	}
	a := aurora["eur"].(float64)
	fmt.Printf("P %s AURORA %11.6f EUR\n", t, a)
	dsdEUR, ok := dsd["eur"].(float64)
	if ok {
		fmt.Printf("P %s DSD %11.6f EUR\n", t, dsdEUR)
	}
	fmt.Printf("P %s ESD %11.6f EUR\n", t, esd["eur"].(float64))
	fmt.Printf("P %s FRAX %11.6f EUR\n", t, frax["eur"].(float64))
	fmt.Printf("P %s FXS %11.6f EUR\n", t, fxs["eur"].(float64))

	stash := os.Getenv("AURORA_STASH")
	if stash != "" {
		ss, err := strconv.ParseFloat(stash, 64)
		if err != nil {
			panic(err)
		}
		total := a * ss
		amount, si := humanize.ComputeSI(total)
		fmt.Fprintf(os.Stderr, "AURORA stash: %6.1f%s EUR (%6.1f BTC)\n",
			amount, si, total/btc)
	}

	total := os.Getenv("AURORA_TOTAL")
	if total != "" {
		ss, err := strconv.ParseFloat(total, 64)
		if err != nil {
			panic(err)
		}
		total := a * ss
		amount, si := humanize.ComputeSI(total)
		fmt.Fprintf(os.Stderr, "AURORA total: %6.1f%s EUR (%6.1f BTC)\n",
			amount, si, total/btc)
	}

	fmt.Fprintf(os.Stderr, "Ethereum/Bitcoin ratio: %.3f\n", eth/btc)
}
