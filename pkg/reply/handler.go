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
		tb.NewKeyboardButton(currency.DollarProviderLabel),
	),
)

type currencyService interface {
	GetLastPrices(providerName string) (*currency.CurrencyPriceList, error)
}

type templateEngine interface {
	ProcessPricesTemplate(priceList *currency.CurrencyPriceList) (string, error)
}

type handler struct {
	tgAPI           *tb.BotAPI
	currencyService currencyService
	templateEngine  templateEngine
	logger          *zap.Logger
}

func NewHandler(b *tb.BotAPI, cs currencyService, t templateEngine, l *zap.Logger) *handler {
	return &handler{
		tgAPI:           b,
		currencyService: cs,
		templateEngine:  t,
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
		case currency.DollarProviderLabel:
			msg.ParseMode = tb.ModeHTML
			msg.Text = h.handleProviderCommand(currency.DollarProviderLabel)
		default:
			msg.Text = "Intenta con /cotizaciones"
		}
	}

	_, err := h.tgAPI.Send(msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to send message for command [%s] to chatID [%d]", update.Message.Text, update.Message.Chat.ID), zap.Error(err))
	}
}

func (h *handler) handleProviderCommand(providerName string) string {
	h.logger.Info(fmt.Sprintf("handle provider command: %s", providerName))
	lastPrices, err := h.currencyService.GetLastPrices(providerName)
	if err != nil {
		h.logger.Error("getting prices", zap.Error(err))
		return errorMsg
	}

	message, err := h.templateEngine.ProcessPricesTemplate(lastPrices)
	if err != nil {
		h.logger.Error("formatting prices template", zap.Error(err))
		return errorMsg
	}

	return message
}
