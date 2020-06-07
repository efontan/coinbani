package provider

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"coinbani/cmd/coinbani/options"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"

	"github.com/pkg/errors"
)

const (
	satoshiResponseExpiration  = 10 * time.Minute
	satoshiARSResponseCacheKey = "satoshi_ars_response"
	satoshiUSDResponseCacheKey = "satoshi_usd_response"
)

var parseSatoshiTResponseFunc = func(r *http.Response) (interface{}, error) {
	var satoshiTResponse satoshiResponse
	err := json.NewDecoder(r.Body).Decode(&satoshiTResponse)
	if err != nil {
		return nil, errors.Wrap(err, "decoding Satoshi response json")
	}
	defer r.Body.Close()

	return satoshiTResponse, nil
}

type satoshiResponse struct {
	Data satoshiData `json:"data"`
}

type satoshiData struct {
	Ticker satoshiTicker `json:"ticker"`
}

type satoshiTicker struct {
	BTC *satoshiPrice `json:"BTC"`
	DAI *satoshiPrice `json:"DAI"`
	ETH *satoshiPrice `json:"ETH"`
}

type satoshiPrice struct {
	BidPrice float64 `json:"bid"`
	AskPrice float64 `json:"ask"`
}

type satoshiTProvider struct {
	config     *options.ProvidersConfig
	restClient client.Http
}

func NewSatoshiTProvider(c *options.ProvidersConfig, r client.Http) *satoshiTProvider {
	return &satoshiTProvider{config: c, restClient: r}
}

func (p *satoshiTProvider) FetchLastPrices() ([]*currency.CurrencyPrice, error) {
	var lastPrices []*currency.CurrencyPrice
	var err error

	// USD
	lastPrices, err = p.fetchPricesForCurrency("ARS", p.config.SatoshiARSURL, satoshiARSResponseCacheKey, lastPrices)
	if err != nil {
		return nil, err
	}

	// USD
	lastPrices, err = p.fetchPricesForCurrency("USD", p.config.SatoshiUSDURL, satoshiUSDResponseCacheKey, lastPrices)
	if err != nil {
		return nil, err
	}

	return lastPrices, nil
}

func (p *satoshiTProvider) fetchPricesForCurrency(currency string, fetchURL string, cacheKey string, lastPrices []*currency.CurrencyPrice) ([]*currency.CurrencyPrice, error) {
	req := &client.GetRequestBuilder{
		Url:             fetchURL,
		CacheKey:        cacheKey,
		CacheExpiration: satoshiResponseExpiration,
		ParseResponse:   parseSatoshiTResponseFunc,
	}

	res, err := p.restClient.Get(req)
	if err != nil {
		return nil, errors.Wrap(err, "fetching prices from SatoshiT service")
	}

	satoshiTResponse := res.(satoshiResponse)

	// DAI
	lastPrices = addCryptocurrencySTPrice(lastPrices, satoshiTResponse.Data.Ticker.DAI, "DAI", currency)
	// BTC
	lastPrices = addCryptocurrencySTPrice(lastPrices, satoshiTResponse.Data.Ticker.BTC, "BTC", currency)
	// ETH
	lastPrices = addCryptocurrencySTPrice(lastPrices, satoshiTResponse.Data.Ticker.ETH, "ETH", currency)

	return lastPrices, nil
}

func addCryptocurrencySTPrice(lastPrices []*currency.CurrencyPrice, price *satoshiPrice, bidCurrency string, askCurrency string) []*currency.CurrencyPrice {
	desc := strings.ToUpper(bidCurrency) + "/" + strings.ToUpper(askCurrency)

	lastPrices = append(lastPrices, &currency.CurrencyPrice{
		Desc:     desc,
		Currency: askCurrency,
		BidPrice: price.BidPrice * 0.99,
		AskPrice: price.AskPrice * 1.01,
	})

	return lastPrices
}
