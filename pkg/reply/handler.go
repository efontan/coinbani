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

var optionsKeyboard = tb.NewReplyKeyboard(
	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton(currency.BBProviderLabel),
		tb.NewKeyboardButton(currency.SatoshiTProviderLabel),
	),
	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton(currency.DolarProviderLabel),
	),
)

type currencyService interface {
	GetLastPrices(providerName string) ([]*currency.CurrencyPriceList, error)
}

type handler struct {
	tgAPI           *tb.BotAPI
	currencyService currencyService
	logger          *zap.Logger
}

func NewHandler(b *tb.BotAPI, cs currencyService, l *zap.Logger) *handler {
	return &handler{
		tgAPI:           b,
		currencyService: cs,
		logger:          l,
	}
}

func (h *handler) HandleReply(update tb.Update) {
	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	h.logger.Info(fmt.Sprintf("handling message [%s] %s", update.Message.From.UserName, update.Message.Text))

	msg := tb.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "cotizaciones":
		msg.Text = "Selecciona una opción para ver las cotizaciones:"
		msg.ReplyMarkup = optionsKeyboard
	default:
		switch update.Message.Text {
		case currency.BBProviderLabel:
			msg.ParseMode = tb.ModeHTML
			msg.Text = h.handleProviderCommand(currency.BBProviderLabel)
		case currency.SatoshiTProviderLabel:
			msg.ParseMode = tb.ModeHTML
			msg.Text = h.handleProviderCommand(currency.SatoshiTProviderLabel)
		case currency.DolarProviderLabel:
			msg.ParseMode = tb.ModeHTML
			msg.Text = h.handleProviderCommand(currency.DolarProviderLabel)
		default:
			msg.Text = "Intenta con /cotizaciones"
		}
	}

	h.tgAPI.Send(msg)
}

func (h *handler) handleProviderCommand(providerName string) string {
	lastPrices, err := h.currencyService.GetLastPrices(providerName)
	if err != nil {
		h.logger.Error("getting prices", zap.Error(err))
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
			if price.PercentChange != "" {
				message = message + fmt.Sprintf("  Variación: %s\n\n", price.PercentChange)
			}
		}
	}

	return message, nil
}
