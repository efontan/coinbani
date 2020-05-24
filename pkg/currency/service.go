package currency

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type currencyProvider interface {
	FetchLastPrices() ([]*CurrencyPrice, error)
}

type service struct {
	bbProvider    currencyProvider
	stProvider    currencyProvider
	dolarProvider currencyProvider
	logger        *zap.Logger
}

func NewService(bb currencyProvider, st currencyProvider, dolar currencyProvider, l *zap.Logger) *service {
	return &service{bbProvider: bb, stProvider: st, dolarProvider: dolar, logger: l}
}

func (s *service) GetLastPrices(providerName string) ([]*CurrencyPriceList, error) {
	switch providerName {
	case BBProviderLabel:
		return s.getBBPrices()
	case SatoshiTProviderLabel:
		return s.getSatoshiTPrices()
	case DolarProviderLabel:
		return s.getDolarPrices()
	default:
		return nil, errors.New("unknown provider")
	}
}

func (s *service) getBBPrices() ([]*CurrencyPriceList, error) {
	var lastPrices []*CurrencyPriceList

	// Fetch BB prices
	bbLastPrices, err := s.bbProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching BB prices")
	}
	lastPrices = append(lastPrices, &CurrencyPriceList{ProviderName: BBProviderLabel, Prices: bbLastPrices})

	return lastPrices, nil
}

func (s *service) getSatoshiTPrices() ([]*CurrencyPriceList, error) {
	var lastPrices []*CurrencyPriceList

	// Fetch SatoshiT prices
	stLastPrices, err := s.stProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching SatoshiT prices")
	}
	lastPrices = append(lastPrices, &CurrencyPriceList{ProviderName: SatoshiTProviderLabel, Prices: stLastPrices})

	return lastPrices, nil
}

func (s *service) getDolarPrices() ([]*CurrencyPriceList, error) {
	var lastPrices []*CurrencyPriceList

	// Fetch Dolar prices
	dolarLastPrices, err := s.dolarProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching Dolar prices")
	}
	lastPrices = append(lastPrices, &CurrencyPriceList{ProviderName: DolarProviderLabel, Prices: dolarLastPrices})

	return lastPrices, nil
}
