package crypto

type CryptocurrencyList struct {
	Exchange string
	Prices   []*CryptocurrencyPrice
}

type CryptocurrencyPrice struct {
	Desc     string
	Currency string
	BidPrice float64
	AskPrice float64
}
