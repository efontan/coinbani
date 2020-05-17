package reply

import (
	"fmt"

	"coinbani/pkg/crypto"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type cryptoService interface {
	GetLastPrices() ([]*crypto.CryptocurrencyList, error)
}

type handler struct {
	bot           *tb.BotAPI
	cryptoService cryptoService
	logger        *zap.Logger
}

func NewHandler(b *tb.BotAPI, l *zap.Logger) *handler {
	return &handler{
		bot:    b,
		logger: l,
	}
}

func (h *handler) HandleReply(update tb.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	h.logger.Info(fmt.Sprintf("[%s] %s", update.Message.From.UserName, update.Message.Text))

	msg := tb.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "help":
		msg.Text = "comandos disponibles: /sayhi - /status - /precios"
	case "sayhi":
		msg.Text = "Hi :)"
	case "status":
		msg.Text = "I'm ok."
	case "precios":
		msg.ParseMode = tb.ModeMarkdown
		msg.Text = h.handlePricesCommand()
	default:
		msg.Text = "Try with /help"
	}

	h.bot.Send(msg)
}

// TODO: implement
func (h *handler) handlePricesCommand() string {
	lastPrices, err := h.crypto.GetLastPrices()
	if err != nil {
		h.logger.Error("getting cryto prices")
	}

	message, err := h.formatPricesMessage(prices)
	if err != nil {
		h.logger.Error("formatting prices message")
	}

	return message
}

func (h *handler) formatPricesMessage() (string, error) {
	return `
#### Buenbit 2.0 ####
	Operación     Compra     Venta
	-------------------------------
#### Satoshi Tango ####
	Operación     Compra     Venta
    -------------------------------
    DAI/ARS           118.5     128
    DAI/USD           134       128
    BTC/ARS           134       128
    BTC/USD           134       128

#### Dolar ####
	Tipo           Compra     Venta
    -------------------------------
    Solidario      $69.7      $65.2
    Blue           $134       $128

*Última Actualización: 16/05/2020 17:45hs*
`
}

var templateMarkdown string = `
{{range $price := .}}
#### {{$exchange}} ####
{{range $element}}
{{.Crypto.Description}}
{{end}}

{{end}}



#### Buenbit 2.0 ####
	Operación     Compra     Venta
	-------------------------------
	{{range $exchange, $price := .}}
	{{$exchange}}  
{{range $element}}{{.Value}}
{{end}}
{{end}}
#### Satoshi Tango ####
	Operación     Compra     Venta
    -------------------------------
    DAI/ARS           118.5     128
    DAI/USD           134       128
    BTC/ARS           134       128
    BTC/USD           134       128

#### Dolar ####
	Tipo           Compra     Venta
    -------------------------------
    Solidario      $69.7      $65.2
    Blue           $134       $128

*Última Actualización: 16/05/2020 17:45hs*
`

// func processTemplate() string {
// 	t := template.New("action")

// 	var err error
// 	t, err = t.ParseFiles("path/to/action.html")
// 	if err != nil {
// 		return err
// 	}

// 	key := "some strings"

// 	data := struct{
// 		Key string
// 	}{
// 		Key: key
// 	}

// 	var tpl bytes.Buffer
// 	if err := t.Execute(&tpl, data); err != nil {
// 		return err
// 	}

// 	result := tpl.String()
// }