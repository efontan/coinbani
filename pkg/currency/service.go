package currency

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type currencyProvider interface {
	FetchLastPrices() ([]*CurrencyPrice, error)
}

type service struct {
	bbProvider     currencyProvider
	stProvider     currencyProvider
	dollarProvider currencyProvider
	logger         *zap.Logger
}

func NewService(bb currencyProvider, st currencyProvider, dollar currencyProvider, l *zap.Logger) *service {
	return &service{bbProvider: bb, stProvider: st, dollarProvider: dollar, logger: l}
}

func (s *service) GetLastPrices(providerName string) (*CurrencyPriceList, error) {
	switch providerName {
	case BBProviderLabel:
		return s.getBBPrices()
	case SatoshiTProviderLabel:
		return s.getSatoshiTPrices()
	case DollarProviderLabel:
		return s.getDollarPrices()
	default:
		return nil, errors.New("unknown provider")
	}
}

func (s *service) getBBPrices() (*CurrencyPriceList, error) {
	// Fetch BB prices
	bbLastPrices, err := s.bbProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching BB prices")
	}

	return &CurrencyPriceList{ProviderName: BBProviderLabel, Prices: bbLastPrices}, nil
}

func (s *service) getSatoshiTPrices() (*CurrencyPriceList, error) {
	// Fetch SatoshiT prices
	stLastPrices, err := s.stProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching SatoshiT prices")
	}

	return &CurrencyPriceList{ProviderName: SatoshiTProviderLabel, Prices: stLastPrices}, nil
}

func (s *service) getDollarPrices() (*CurrencyPriceList, error) {
	// Fetch Dollar prices
	dollarLastPrices, err := s.dollarProvider.FetchLastPrices()
	if err != nil {
		return nil, errors.Wrap(err, "fetching Dollar prices")
	}

	return &CurrencyPriceList{ProviderName: DollarProviderLabel, Prices: dollarLastPrices}, nil
}
