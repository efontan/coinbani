package crypto

import (
	"go.uber.org/zap"
)


type exchange interface {
	FetchLastPrices() ([]*CryptocurrencyPrice, error)
}

type service struct {
	bbExchange exchange
	logger     *zap.Logger
}

func NewService(bb exchange, l *zap.Logger) *service {
	return &service{bbExchange: bb, logger: l}
}

func (s *service) GetLastPrices() ([]*CryptocurrencyList, error) {
	var lastPrices []*CryptocurrencyList

	// Fetch BB prices
	bbLastPrices, err := s.bbExchange.FetchLastPrices()
	if err != nil {
		return nil, err
	}
	lastPrices = append(lastPrices, &CryptocurrencyList{Exchange: "Buenbit 2.0", Prices: bbLastPrices})

	// Fetch Satoshi Tango prices
	// TODO: implement

	return lastPrices, nil
}
