package template

import (
	"fmt"

	"coinbani/pkg/currency"

	"github.com/alexeyco/simpletable"
)

type priceData struct {
	ProviderName string
	PricesTable  string
}

type tableFormatter interface {
	FormatPricesTable(prices []*currency.CurrencyPrice) (content string, err error)
}

type simpleTableFormatter struct {
}

func NewTableFormatter() *simpleTableFormatter {
	return &simpleTableFormatter{}
}

func (t *simpleTableFormatter) FormatPricesTable(prices []*currency.CurrencyPrice) (content string, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("unknown panic formatting table data: %v", r)
			}
			content = ""
		}
	}()

	table := simpletable.New()
	addPercentage := false

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: "Moneda"},
			{Align: simpletable.AlignLeft, Text: "Compra"},
			{Align: simpletable.AlignLeft, Text: "Venta"},
		},
	}

	for _, price := range prices {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: price.Desc},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%.2f", price.BidPrice)},
			{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%.2f", price.AskPrice)},
		}

		if price.PercentChange != "" {
			addPercentage = true
			r = append(r, &simpletable.Cell{Align: simpletable.AlignLeft, Text: price.PercentChange})
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	if addPercentage {
		table.Header.Cells = append(table.Header.Cells, &simpletable.Cell{Align: simpletable.AlignLeft, Text: "%"})
	}

	table.SetStyle(simpletable.StyleCompactLite)
	return table.String(), nil
}
