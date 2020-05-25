package template

import (
	"bytes"
	"text/template"

	"coinbani/pkg/currency"

	"github.com/pkg/errors"
)

const (
	PricesTemplate = `
<strong>{{.ProviderName}}</strong>

<pre>
{{.PricesTable}}
</pre>
`
)

type templateEngine struct {
	tableFormatter tableFormatter
}

func NewEngine() *templateEngine {
	return &templateEngine{tableFormatter: NewTableFormatter()}
}

func (e *templateEngine) FormatPricesMessage(priceList *currency.CurrencyPriceList) (string, error) {
	pricesTable, err := e.tableFormatter.FormatPricesTable(priceList.Prices)
	if err != nil {
		return "", errors.Wrap(err, "formatting prices table")
	}

	data := &priceData{
		ProviderName: priceList.ProviderName,
		PricesTable:  pricesTable,
	}
	return e.processTemplate(PricesTemplate, data)
}

func (e *templateEngine) processTemplate(t string, data interface{}) (string, error) {
	tmpl, err := template.New("tmpl").Parse(t)
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
