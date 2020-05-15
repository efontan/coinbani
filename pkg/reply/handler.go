package reply

import (
	"fmt"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type handler struct {
	bot    *tb.BotAPI
	logger *zap.Logger
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
		msg.Text = "type /sayhi or /status."
	case "sayhi":
		msg.Text = "Hi :)"
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = "Try with /help"
	}

	h.bot.Send(msg)
}
