package currency

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type provider interface {
	FetchLastPrices() ([]*CurrencyPrice, error)
}

type service struct {
	bbProvider    provider
	stProvider    provider
	dolarProvider provider
	logger        *zap.Logger
}

func NewService(bb provider, st provider, dolar provider, l *zap.Logger) *service {
	return &service{bbProvider: bb, stProvider: st, dolarProvider: dolar, logger: l}
}

func (s *service) GetLastPrices() ([]*CurrencyPriceList, error) {
	var lastPrices []*CurrencyPriceList

	// Fetch BB prices
	bbLastPrices, err := s.bbProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching BB prices")
	}
	lastPrices = append(lastPrices, &CurrencyPriceList{ProviderName: "Buenbit 2.0", Prices: bbLastPrices})

	// Fetch SatoshiT prices
	stLastPrices, err := s.stProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching SatoshiT prices")
	}
	lastPrices = append(lastPrices, &CurrencyPriceList{ProviderName: "Satoshi Tango", Prices: stLastPrices})

	// Fetch Dolar prices
	dolarLastPrices, err := s.dolarProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching Dolar prices")
	}
	lastPrices = append(lastPrices, &CurrencyPriceList{ProviderName: "Dolar", Prices: dolarLastPrices})

	return lastPrices, nil
}
