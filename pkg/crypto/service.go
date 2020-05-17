package crypto

import "go.uber.org/zap"

type service struct {
	logger *zap.Logger
}

func NewService(l *zap.Logger) *service {
	return &service{logger: l}
}

// TODO: implement
func (s *service) GetLastPrices() ([]*CryptocurrencyList, error) {
	var lastPrices []*CryptocurrencyList

	// Fetch BuenBit prices

	// Fetch Satoshi Tango prices

	return lastPrices, nil
}

func (s *service) FetchBuenBitLastPrices() ([]*CryptocurrencyPrice, error) {
	return nil, nil
}
