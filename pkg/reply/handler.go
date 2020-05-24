package reply

import (
	"fmt"

	"coinbani/pkg/currency"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	errorMsg = "Lo sentimos, ha ocurrido un error intenta más tarde"
)

var testKeyboard = tb.NewReplyKeyboard(
	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton("BuenBit"),
		tb.NewKeyboardButton("Satoshi Tango"),
	),
	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton("Dolar"),
	),
)

type cryptoService interface {
	GetLastPrices() ([]*currency.CurrencyPriceList, error)
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
		msg.Text = "Comandos disponibles:\n /cotizaciones"
	case "test":
		msg.Text = "Selecciona una opción:"
		msg.ReplyMarkup = testKeyboard
	case "cotizaciones":
		msg.ParseMode = tb.ModeHTML
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
		h.logger.Error("getting cryto prices", zap.Error(err))
		return errorMsg
	}

	message, err := h.formatPricesMessage(lastPrices)
	if err != nil {
		h.logger.Error("formatting prices message", zap.Error(err))
		return errorMsg
	}

	return message
}

func (h *handler) formatPricesMessage(lastPrices []*currency.CurrencyPriceList) (string, error) {
	message := ""

	for _, p := range lastPrices {
		message = message + fmt.Sprint("-------------------------------------\n")
		message = message + fmt.Sprintf("<strong>%s\n</strong>", p.ProviderName)
		message = message + fmt.Sprint("-------------------------------------\n")

		for _, price := range p.Prices {
			message = message + fmt.Sprintf("%s\n", price.Desc)
			message = message + fmt.Sprintf("  Compra: %.2f\n", price.BidPrice)
			message = message + fmt.Sprintf("  Venta: %.2f\n\n", price.AskPrice)
		}
	}

	return message, nil
}
