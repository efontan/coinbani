package main

import (
	"fmt"
	golog "log"
	"net/http"

	"coinbani/cmd/coinbani/options"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"
	"coinbani/pkg/currency/provider"
	"coinbani/pkg/reply"
	"coinbani/pkg/template"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func main() {
	cfg := options.NewConfig()
	bot, err := tb.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		golog.Panic(err)
	}

	logger, err := options.GetLogger(cfg.Log)
	defer logger.Sync()
	logger.Debug("initializing coinbani bot whit config", zap.String("cfg", fmt.Sprintf("%+v", cfg)))

	// Setup telegram bot
	bot.Debug = cfg.Bot.Debug
	logger.Info(fmt.Sprintf("authorized on account %s", bot.Self.UserName))

	logger.Info("setting up webhook", zap.String("CallbackURL", cfg.Application.CallbackURL))
	_, err = bot.SetWebhook(tb.NewWebhook(cfg.Application.CallbackURL + bot.Token))
	if err != nil {
		logger.Fatal("setting up webhook", zap.Error(err))
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		logger.Fatal("getting webhook info", zap.Error(err))
	}
	if info.LastErrorDate != 0 {
		logger.Error(fmt.Sprintf("Telegram callback failed: %s", info.LastErrorMessage))
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(":"+cfg.Application.Port, nil)

	// setup services
	restClient := client.NewRestClient()
	bbProvider := provider.NewBBProvider(cfg.Providers, restClient)
	satoshiTProvider := provider.NewSatoshiTProvider(cfg.Providers, restClient)
	dollarProvider := provider.NewDollarProvider(cfg.Providers, restClient)

	currencyService := currency.NewService(bbProvider, satoshiTProvider, dollarProvider, logger)
	templateEngine := template.NewEngine()
	replyHandler := reply.NewHandler(bot, currencyService, templateEngine, logger)

	logger.Info("coinbani bot successfully started!")

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				logger.Info("channel closed, exiting")
				break
			}

			go replyHandler.HandleReply(update)
		}
	}
}
