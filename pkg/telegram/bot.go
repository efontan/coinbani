package telegram

import (
	"coinbani/cmd/coinbani/options"
	"fmt"
	"net/http"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Bot interface {
	Send(c tb.Chattable) (tb.Message, error)
}

type tbbot struct {
	bot    *tb.BotAPI
	config *options.Config
	logger *zap.Logger
}

func NewBot(c *options.Config, l *zap.Logger) (*tbbot, error) {
	b, err := tb.NewBotAPI(c.Bot.Token)
	if err != nil {
		return nil, errors.Wrap(err, "initializing Telegram bot")
	}
	b.Debug = c.Bot.Debug
	l.Info(fmt.Sprintf("authorized on account %s", b.Self.UserName))

	return &tbbot{
		bot:    b,
		logger: l,
	}, nil
}

func (b *tbbot) GetUpdatesChan() (tb.UpdatesChannel, error) {
	u := tb.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *tbbot) GetUpdatesChanForWebhook(pattern string) (tb.UpdatesChannel, error) {
	b.logger.Info("setting up webhook", zap.String("CallbackURL", b.config.Application.CallbackURL))
	_, err := b.bot.SetWebhook(tb.NewWebhook(b.config.Application.CallbackURL + b.bot.Token))
	if err != nil {
		return nil, errors.Wrap(err, "setting up bot webhook")
	}
	info, err := b.bot.GetWebhookInfo()
	if err != nil {
		return nil, errors.Wrap(err, "getting bot webhook info")
	}
	if info.LastErrorDate != 0 {
		return nil, errors.Wrapf(err, "Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := b.bot.ListenForWebhook("/" + b.bot.Token)
	go http.ListenAndServe(":"+b.config.Application.Port, nil)

	return updates, nil
}

func (b *tbbot) Send(c tb.Chattable) (tb.Message, error) {
	return b.bot.Send(c)
}
