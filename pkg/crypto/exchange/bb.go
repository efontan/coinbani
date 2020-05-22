package exchange

import (
	"coinbani/cmd/options"
	"coinbani/pkg/cache"
	"coinbani/pkg/crypto"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	Expiration         = 30 * time.Minute
	BBResponseCacheKey = "bb_response"
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

type bbExchange struct {
	httpClient *http.Client
	config     *options.ExchangeConfig
	cache      cache.Cache
}

func NewBBExchange(c *options.ExchangeConfig, httpClient *http.Client, cache cache.Cache) *bbExchange {
	return &bbExchange{config: c, httpClient: httpClient, cache: cache}
}

func (e *bbExchange) FetchLastPrices() ([]*crypto.CryptocurrencyPrice, error) {
	var lastPrices []*crypto.CryptocurrencyPrice
	r, err := e.httpClient.Get(e.config.BBURL)
	if err != nil {
		return nil, errors.Wrap(err, "fetching prices from BB service")
	}
	defer r.Body.Close()

	var bbResponse BBResponse

	cachedResponse, found := e.cache.Get(BBResponseCacheKey)
	if !found {
		// fetch from service
		err = json.NewDecoder(r.Body).Decode(&bbResponse)
		if err != nil || bbResponse.Object == nil {
			return nil, errors.Wrap(err, "decoding BB response json")
		}
		e.cache.Set(BBResponseCacheKey, bbResponse, Expiration)
	} else {
		// fetch from cache
		bbResponse = cachedResponse.(BBResponse)
	}

	// DAI ARS
	lastPrices = addCryptocurrencyPrice(lastPrices, bbResponse.Object.DaiARS)
	// DAI USD
	lastPrices = addCryptocurrencyPrice(lastPrices, bbResponse.Object.DaiUSD)
	// DAI USD
	lastPrices = addCryptocurrencyPrice(lastPrices, bbResponse.Object.BTCARS)

	return lastPrices, nil
}

func addCryptocurrencyPrice(lastPrices []*crypto.CryptocurrencyPrice, price *BBPrice) []*crypto.CryptocurrencyPrice {
	desc := strings.ToUpper(price.BidCurrency) + "/" + strings.ToUpper(price.AskCurrency)

	lastPrices = append(lastPrices, &crypto.CryptocurrencyPrice{
		Desc:     desc,
		Currency: strings.Replace(price.Currency, "$", "S", -1),
		BidPrice: price.BidPrice,
		AskPrice: price.AskPrice,
	})

	return lastPrices
}
