package provider

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"coinbani/cmd/coinbani/options"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"

	"github.com/pkg/errors"
)

const (
	dollarResponseExpiration = 30 * time.Minute
	dollarResponseCacheKey   = "dollar_response"
	dollarOfficial           = "Dolar Oficial"
	dollarBlue               = "Dolar Blue"
	dollarMEP                = "Dolar Bolsa"
	dollarCCL                = "Dolar Contado con Liqui"
	dollarSaving             = "Ahorro"
)

var namesMap = map[string]string{
	"Bolsa":             "MEP",
	"Contado con Liqui": "CCL",
}

var parseDollarResponseFunc = func(r *http.Response) (interface{}, error) {
	var dollarResponse *dollarRateResponse
	err := json.NewDecoder(r.Body).Decode(&dollarResponse)
	if err != nil || len(*dollarResponse) < 2 {
		return nil, errors.Wrap(err, "decoding dollar response json")
	}
	defer r.Body.Close()

	return dollarResponse, nil
}

type dollarRateResponse []struct {
	Price dollarPrice `json:"casa"`
}

type dollarPrice struct {
	BidPrice      string `json:"compra"`
	AskPrice      string `json:"venta"`
	Name          string `json:"nombre"`
	PercentChange string `json:"variacion"`
}

type dollarProvider struct {
	config     *options.ProvidersConfig
	restClient client.Http
}

func NewDollarProvider(c *options.ProvidersConfig, r client.Http) *dollarProvider {
	return &dollarProvider{config: c, restClient: r}
}

func (d *dollarProvider) FetchLastPrices() ([]*currency.CurrencyPrice, error) {
	req := &client.GetRequestBuilder{
		Url:             d.config.DollarURL,
		CacheKey:        dollarResponseCacheKey,
		CacheExpiration: dollarResponseExpiration,
		ParseResponse:   parseDollarResponseFunc,
	}

	res, err := d.restClient.Get(req)
	if err != nil {
		return nil, errors.Wrap(err, "fetching prices from dollar service")
	}

	dollarResponse := res.(*dollarRateResponse)

	prices := filterPrices(*dollarResponse)
	prices, err = d.addDollarSaving(prices)
	if err != nil {
		return nil, errors.Wrap(err, "addding dollar saving")
	}

	var lastPrices []*currency.CurrencyPrice
	for _, p := range prices {
		lastPrices = addDollarPrices(lastPrices, p)
	}

	return lastPrices, nil
}

func filterPrices(response dollarRateResponse) []dollarPrice {
	prices := make([]dollarPrice, 0)
	for _, p := range response {
		if p.Price.Name == dollarOfficial || p.Price.Name == dollarBlue || p.Price.Name == dollarMEP || p.Price.Name == dollarCCL {
			prices = append(prices, p.Price)
		}
	}
	return prices
}

func (d *dollarProvider) addDollarSaving(prices []dollarPrice) ([]dollarPrice, error) {
	var official *dollarPrice
	for _, p := range prices {
		if p.Name == dollarOfficial {
			official = &p
			break
		}
	}

	if official == nil {
		return nil, errors.New("official dollar not found in list")
	}

	b := strings.Replace(official.BidPrice, ",", ".", 1)
	bidPrice, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return nil, errors.New("error parsing official dollar bid price")
	}

	a := strings.Replace(official.AskPrice, ",", ".", 1)
	askPrice, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return nil, errors.New("error parsing official dollar ask price")
	}

	savingDollar := dollarPrice{
		Name:          dollarSaving,
		BidPrice:      strconv.FormatFloat(bidPrice*d.config.DollarSavingTax, 'f', 2, 64),
		AskPrice:      strconv.FormatFloat(askPrice*d.config.DollarSavingTax, 'f', 2, 64),
		PercentChange: official.PercentChange,
	}

	prices = append(prices, savingDollar)
	return prices, nil
}

func addDollarPrices(lastPrices []*currency.CurrencyPrice, price dollarPrice) []*currency.CurrencyPrice {
	bidStr := replaceComa(price.BidPrice)
	bidPrice, err := strconv.ParseFloat(bidStr, 32)
	if err != nil {
		return lastPrices
	}

	askStr := replaceComa(price.AskPrice)
	askPrice, err := strconv.ParseFloat(askStr, 32)
	if err != nil {
		return lastPrices
	}

	lastPrices = append(lastPrices, &currency.CurrencyPrice{
		Desc:          formatDollarName(price.Name),
		Currency:      "USD",
		BidPrice:      math.Round(bidPrice*100) / 100,
		AskPrice:      math.Round(askPrice*100) / 100,
		PercentChange: formatPercent(price.PercentChange),
	})

	return lastPrices
}

func formatDollarName(v string) string {
	name := strings.Replace(v, "Dolar ", "", -1)

	v, ok := namesMap[name]
	if !ok {
		return name
	}

	return v
}

func replaceComa(value string) string {
	return strings.Replace(value, ",", ".", -1)
}

func formatPercent(percentChange string) string {
	if !strings.Contains(percentChange, "-") {
		return "+" + percentChange
	}
	return percentChange
}
