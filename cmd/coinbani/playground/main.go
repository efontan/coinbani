package main

import (
	"fmt"
	golog "log"

	"coinbani/cmd/coinbani/options"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"
	"coinbani/pkg/currency/provider"
	"coinbani/pkg/reply"
	"coinbani/pkg/template"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

// For local testing purposes
func main() {
	cfg := options.NewConfig()
	bot, err := tb.NewBotAPI(cfg.Bot.TokenBeta)
	if err != nil {
		golog.Panic(err)
	}

	logger, err := options.GetLogger(cfg.Log)
	defer logger.Sync()
	logger.Debug("initializing coinbani bot whit config", zap.String("cfg", fmt.Sprintf("%+v", cfg)))

	// Setup telegram bot
	bot.Debug = cfg.Bot.Debug
	logger.Info(fmt.Sprintf("authorized on account %s", bot.Self.UserName))

	u := tb.NewUpdate(0)
	u.Timeout = 60
	logger.Info("starting channel for getting bot updates")
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.Fatal("starting bot channel", zap.Error(err))
	}

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
