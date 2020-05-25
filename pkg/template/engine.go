package template

import (
	"bytes"
	"text/template"

	"coinbani/cmd/options"
	"coinbani/pkg/currency"

	"github.com/pkg/errors"
)

const (
	pricesSummaryTemplateFileName = "prices.gohtml"
)

type templateEngine struct {
	config         *options.TemplateConfig
	tableFormatter tableFormatter
}

func NewEngine(c *options.TemplateConfig) *templateEngine {
	return &templateEngine{config: c, tableFormatter: NewTableFormatter()}
}

func (e *templateEngine) ProcessPricesTemplate(priceList *currency.CurrencyPriceList) (string, error) {
	pricesTable, err := e.tableFormatter.FormatPricesTable(priceList.Prices)
	if err != nil {
		return "", errors.Wrap(err, "formating prices table")
	}

	t := &priceTemplate{
		ProviderName: priceList.ProviderName,
		PricesTable:  pricesTable,
	}
	return e.processTemplate(e.config.TemplatesDir+pricesSummaryTemplateFileName, t)
}

func (e *templateEngine) processTemplate(fileName string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		return "", errors.Wrap(err, "parsing template file")
	}

	return process(tmpl, data)
}

func process(tmpl *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
		return "", errors.Wrap(err, "processing template")
	}

	return buf.String(), nil
}
