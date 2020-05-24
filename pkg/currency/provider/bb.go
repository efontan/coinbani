package provider

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"coinbani/cmd/options"
	"coinbani/pkg/cache"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"

	"github.com/pkg/errors"
)

const (
	BBResponseExpiration = 20 * time.Minute
	BBResponseCacheKey   = "bb_response"
)

type BBResponse struct {
	Object *BBObject `json:"object"`
}

type BBObject struct {
	DaiARS *BBPrice `json:"daiars"`
	DaiUSD *BBPrice `json:"daiusd"`
	BTCARS *BBPrice `json:"btcars"`
}

type BBPrice struct {
	BidPrice           float64 `json:"purchase_price,string"`
	BidCurrency        string  `json:"bid_currency"`
	AskPrice           float64 `json:"selling_price,string"`
	AskCurrency        string  `json:"ask_currency"`
	PriceChangePercent string  `json:"price_change_percent"`
	Currency           string  `json:"currency"`
	MarketIdentifier   string  `json:"market_identifier"`
}

type bbProvider struct {
	httpClient client.Http
	config     *options.ProvidersConfig
	cache      cache.Cache
}

func NewBBProvider(c *options.ProvidersConfig, httpClient client.Http, cache cache.Cache) *bbProvider {
	return &bbProvider{config: c, httpClient: httpClient, cache: cache}
}

func (e *bbProvider) FetchLastPrices() ([]*currency.CurrencyPrice, error) {
	var lastPrices []*currency.CurrencyPrice

	var bbResponse BBResponse
	cachedResponse, found := e.cache.Get(BBResponseCacheKey)

	if !found {
		// fetch from service
		r, err := e.httpClient.Get(e.config.BBURL)
		if err != nil {
			return nil, errors.Wrap(err, "fetching prices from BB service")
		}
		defer r.Body.Close()

		if r.StatusCode != http.StatusOK {
			return nil, errors.Wrap(err, "fetching prices from BB service")
		}

		err = json.NewDecoder(r.Body).Decode(&bbResponse)
		if err != nil || bbResponse.Object == nil {
			return nil, errors.Wrap(err, "decoding BB response json")
		}
		e.cache.Set(BBResponseCacheKey, bbResponse, BBResponseExpiration)
	} else {
		// fetch from cache
		bbResponse = cachedResponse.(BBResponse)
	}

	// DAI ARS
	lastPrices = addCryptocurrencyBBPrice(lastPrices, bbResponse.Object.DaiARS)
	// DAI USD
	lastPrices = addCryptocurrencyBBPrice(lastPrices, bbResponse.Object.DaiUSD)
	// BTC ARS
	lastPrices = addCryptocurrencyBBPrice(lastPrices, bbResponse.Object.BTCARS)

	return lastPrices, nil
}

func addCryptocurrencyBBPrice(lastPrices []*currency.CurrencyPrice, price *BBPrice) []*currency.CurrencyPrice {
	desc := strings.ToUpper(price.BidCurrency) + "/" + strings.ToUpper(price.AskCurrency)

	lastPrices = append(lastPrices, &currency.CurrencyPrice{
		Desc:          desc,
		Currency:      strings.Replace(price.Currency, "$", "S", -1),
		BidPrice:      price.BidPrice,
		AskPrice:      price.AskPrice,
		PercentChange: price.PriceChangePercent,
	})

	return lastPrices
}
