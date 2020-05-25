package provider

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"coinbani/cmd/options"
	"coinbani/pkg/cache"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"

	"github.com/pkg/errors"
)

const (
	DollarResponseExpiration = 60 * time.Minute
	DollarResponseCacheKey   = "dollar_response"
)

var namesMap = map[string]string{
	"Bolsa":             "MEP",
	"Contado con Liqui": "CCL",
}

type DollarRateResponse []struct {
	Price DollarPrice `json:"casa"`
}

type DollarPrice struct {
	BidPrice      string `json:"compra"`
	AskPrice      string `json:"venta"`
	Name          string `json:"nombre"`
	PercentChange string `json:"variacion"`
}

type dollarProvider struct {
	config     *options.ProvidersConfig
	httpClient client.Http
	cache      cache.Cache
}

func NewDollarProvider(c *options.ProvidersConfig, httpClient client.Http, cache cache.Cache) *dollarProvider {
	return &dollarProvider{config: c, httpClient: httpClient, cache: cache}
}

func (e *dollarProvider) FetchLastPrices() ([]*currency.CurrencyPrice, error) {
	var lastPrices []*currency.CurrencyPrice

	var dollarResponse DollarRateResponse
	cachedResponse, found := e.cache.Get(DollarResponseCacheKey)

	if !found {
		// fetch from service
		r, err := e.httpClient.Get(e.config.DollarURL)
		if err != nil {
			return nil, errors.Wrap(err, "fetching prices from dollar service")
		}
		defer r.Body.Close()

		if r.StatusCode != http.StatusOK {
			return nil, errors.Wrap(err, "fetching prices from dollar service")
		}

		err = json.NewDecoder(r.Body).Decode(&dollarResponse)
		if err != nil || len(dollarResponse) < 2 {
			return nil, errors.Wrap(err, "decoding dollar response json")
		}
		e.cache.Set(DollarResponseCacheKey, dollarResponse, DollarResponseExpiration)
	} else {
		// fetch from cache
		dollarResponse = cachedResponse.(DollarRateResponse)
	}

	prices := filterPRices(dollarResponse)
	for _, p := range prices {
		lastPrices = addDollarPrices(lastPrices, p)
	}

	return lastPrices, nil
}

func filterPRices(response DollarRateResponse) []DollarPrice {
	prices := make([]DollarPrice, 0)
	for _, p := range response {
		if p.Price.Name == "Dolar Oficial" || p.Price.Name == "Dolar Blue" || p.Price.Name == "Dolar Bolsa" || p.Price.Name == "Dolar Contado con Liqui" {
			prices = append(prices, p.Price)
		}
	}
	return prices
}

func addDollarPrices(lastPrices []*currency.CurrencyPrice, price DollarPrice) []*currency.CurrencyPrice {
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
		BidPrice:      bidPrice,
		AskPrice:      askPrice,
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
