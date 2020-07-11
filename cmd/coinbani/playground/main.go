package main

import (
	"coinbani/pkg/telegram"
	"fmt"

	"coinbani/cmd/coinbani/options"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"
	"coinbani/pkg/currency/provider"
	"coinbani/pkg/reply"
	"coinbani/pkg/template"

	"go.uber.org/zap"
)

// For local testing purposes
func main() {
	cfg := options.NewConfig()
	logger, err := options.GetLogger(cfg.Log)
	defer logger.Sync()
	logger.Debug("initializing coinbani bot whit config", zap.String("cfg", fmt.Sprintf("%+v", cfg)))

	bot, err := telegram.NewBot(cfg, logger)
	if err != nil {
		logger.Fatal("initializing telegram bot", zap.Error(err))
	}
	logger.Info("starting channel for getting bot updates")
	updates, err := bot.GetUpdatesChan()
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
