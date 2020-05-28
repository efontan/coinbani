package provider

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"coinbani/cmd/coinbani/options"
	"coinbani/pkg/cache"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"

	"github.com/pkg/errors"
)

const (
	SatoshiResponseExpiration  = 20 * time.Minute
	SatoshiARSResponseCacheKey = "satoshi_ars_response"
	SatoshiUSDResponseCacheKey = "satoshi_usd_response"
)

type SatoshiResponse struct {
	Data *SatoshiData `json:"data"`
}

type SatoshiData struct {
	Ticker *SatoshiTicker `json:"ticker"`
}

type SatoshiTicker struct {
	BTC *SatoshiPrice `json:"BTC"`
	DAI *SatoshiPrice `json:"DAI"`
	ETH *SatoshiPrice `json:"ETH"`
}

type SatoshiPrice struct {
	BidPrice float64 `json:"bid"`
	AskPrice float64 `json:"ask"`
}

type satoshiTProvider struct {
	config     *options.ProvidersConfig
	httpClient client.Http
	cache      cache.Cache
}

func NewSatoshiTProvider(c *options.ProvidersConfig, httpClient client.Http, cache cache.Cache) *satoshiTProvider {
	return &satoshiTProvider{config: c, httpClient: httpClient, cache: cache}
}

func (e *satoshiTProvider) FetchLastPrices() ([]*currency.CurrencyPrice, error) {
	var lastPrices []*currency.CurrencyPrice
	var err error

	// USD
	lastPrices, err = e.fetchPricesForCurrency("ARS", e.config.SatoshiARSURL, SatoshiARSResponseCacheKey, lastPrices)
	if err != nil {
		return nil, err
	}

	// USD
	lastPrices, err = e.fetchPricesForCurrency("USD", e.config.SatoshiUSDURL, SatoshiUSDResponseCacheKey, lastPrices)
	if err != nil {
		return nil, err
	}

	return lastPrices, nil
}

func (e *satoshiTProvider) fetchPricesForCurrency(currency string, fetchURL string, cacheKey string, lastPrices []*currency.CurrencyPrice) ([]*currency.CurrencyPrice, error) {
	var satoshiTResponse SatoshiResponse
	cachedResponse, found := e.cache.Get(cacheKey)

	if !found {
		// fetch from service
		r, err := e.httpClient.Get(fetchURL)
		if err != nil {
			return nil, errors.Wrap(err, "fetching prices from SatoshiT service")
		}
		defer r.Body.Close()

		if r.StatusCode != http.StatusOK {
			return nil, errors.Wrap(err, "fetching prices from SatoshiT service")
		}

		err = json.NewDecoder(r.Body).Decode(&satoshiTResponse)
		if err != nil || satoshiTResponse.Data == nil {
			return nil, errors.Wrap(err, "decoding Satoshi response json")
		}
		e.cache.Set(cacheKey, satoshiTResponse, SatoshiResponseExpiration)
	} else {
		// fetch from cache
		satoshiTResponse = cachedResponse.(SatoshiResponse)
	}

	// DAI
	lastPrices = addCryptocurrencySTPrice(lastPrices, satoshiTResponse.Data.Ticker.DAI, "DAI", currency)
	// BTC
	lastPrices = addCryptocurrencySTPrice(lastPrices, satoshiTResponse.Data.Ticker.BTC, "BTC", currency)
	// ETH
	lastPrices = addCryptocurrencySTPrice(lastPrices, satoshiTResponse.Data.Ticker.ETH, "ETH", currency)

	return lastPrices, nil
}

func addCryptocurrencySTPrice(lastPrices []*currency.CurrencyPrice, price *SatoshiPrice, bidCurrency string, askCurrency string) []*currency.CurrencyPrice {
	desc := strings.ToUpper(bidCurrency) + "/" + strings.ToUpper(askCurrency)

	lastPrices = append(lastPrices, &currency.CurrencyPrice{
		Desc:     desc,
		Currency: askCurrency,
		BidPrice: price.BidPrice,
		AskPrice: price.AskPrice,
	})

	return lastPrices
}
