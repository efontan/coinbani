package reply

import (
	"fmt"
	"time"

	"coinbani/pkg/crypto"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	errorMsg   = "Lo sentimos, ha ocurrido un error intenta más tarde"
)

type cryptoService interface {
	GetLastPrices() ([]*crypto.CryptocurrencyList, error)
}

type handler struct {
	tgAPI         *tb.BotAPI
	cryptoService cryptoService
	logger        *zap.Logger
}

func NewHandler(b *tb.BotAPI, cs cryptoService, l *zap.Logger) *handler {
	return &handler{
		tgAPI:         b,
		cryptoService: cs,
		logger:        l,
	}
}

func (h *handler) HandleReply(update tb.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	h.logger.Info(fmt.Sprintf("handling message [%s] %s", update.Message.From.UserName, update.Message.Text))

	msg := tb.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "help":
		msg.Text = "Comandos disponibles:\n /sayhi - /status - /cotizaciones"
	case "sayhi":
		msg.Text = "Hi :)"
	case "status":
		msg.Text = "I'm ok."
	case "cotizaciones":
		msg.ParseMode = "markdown"
		msg.Text = h.handlePricesCommand()
	default:
		msg.Text = "Intenta con /help"
	}

	h.tgAPI.Send(msg)
}

// TODO: implement
func (h *handler) handlePricesCommand() string {
	lastPrices, err := h.cryptoService.GetLastPrices()
	if err != nil {
		h.logger.Error("getting cryto prices")
		return errorMsg
	}

	message, err := h.formatPricesMessage(lastPrices)
	if err != nil {
		h.logger.Error("formatting prices message")
		return errorMsg
	}

	return message
}

func (h *handler) formatPricesMessage(lastPrices []*crypto.CryptocurrencyList) (string, error) {
	return `
#### Buenbit 2.0 ####
	Operación     Compra     Venta
	-------------------------------
	DAI/ARS        134.5    138
    DAI/USD        1.03     1.07

#### Satoshi Tango ####
	Operación     Compra     Venta
    -------------------------------
    DAI/ARS        134.5    138
    DAI/USD        1.03     1.07
    BTC/ARS        134.5    138
    BTC/USD        134.5    138

#### Dolar ####
	Tipo           Compra     Venta
    -------------------------------
    Solidario      $69.7      $65.2
    Blue           $134       $128

*Última Actualización: 16/05/2020 17:45hs*
`, nil
}
