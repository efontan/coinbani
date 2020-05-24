package currency

type CurrencyPriceList struct {
	ProviderName string
	Prices       []*CurrencyPrice
}

type CurrencyPrice struct {
	Desc          string
	Currency      string
	BidPrice      float64
	AskPrice      float64
	PercentChange string
}
