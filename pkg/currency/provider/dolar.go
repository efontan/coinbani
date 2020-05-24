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
	DolarResponseExpiration = 60 * time.Minute
	DolarResponseCacheKey   = "dolar_response"
)

type DolarRateResponse []struct {
	Price DolarPrice `json:"casa"`
}

type DolarPrice struct {
	BidPrice      string `json:"compra"`
	AskPrice      string `json:"venta"`
	Name          string `json:"nombre"`
	PercentChange string `json:"variacion"`
}

type dolarProvider struct {
	config     *options.ProvidersConfig
	httpClient client.Http
	cache      cache.Cache
}

func NewDolarProvider(c *options.ProvidersConfig, httpClient client.Http, cache cache.Cache) *dolarProvider {
	return &dolarProvider{config: c, httpClient: httpClient, cache: cache}
}

func (e *dolarProvider) FetchLastPrices() ([]*currency.CurrencyPrice, error) {
	var lastPrices []*currency.CurrencyPrice

	var dolarResponse DolarRateResponse
	cachedResponse, found := e.cache.Get(DolarResponseCacheKey)

	if !found {
		// fetch from service
		r, err := e.httpClient.Get(e.config.DolarURL)
		if err != nil {
			return nil, errors.Wrap(err, "fetching prices from dolar service")
		}
		defer r.Body.Close()

		if r.StatusCode != http.StatusOK {
			return nil, errors.Wrap(err, "fetching prices from dolar service")
		}

		err = json.NewDecoder(r.Body).Decode(&dolarResponse)
		if err != nil || len(dolarResponse) < 2 {
			return nil, errors.Wrap(err, "decoding dolar response json")
		}
		e.cache.Set(DolarResponseCacheKey, dolarResponse, DolarResponseExpiration)
	} else {
		// fetch from cache
		dolarResponse = cachedResponse.(DolarRateResponse)
	}

	lastPrices = addDolarPrices(lastPrices, dolarResponse[0].Price)
	lastPrices = addDolarPrices(lastPrices, dolarResponse[1].Price)

	return lastPrices, nil
}

func addDolarPrices(lastPrices []*currency.CurrencyPrice, price DolarPrice) []*currency.CurrencyPrice {
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
		Desc:          price.Name,
		Currency:      "USD",
		BidPrice:      bidPrice,
		AskPrice:      askPrice,
		PercentChange: formatPercent(price.PercentChange),
	})

	return lastPrices
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
