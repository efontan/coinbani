package main

import (
	"context"
	"fmt"
	"log"

	"coinbani/cmd/coinbani/options"
	"coinbani/pkg/client"
	"coinbani/pkg/currency"
	"coinbani/pkg/currency/provider"
	"coinbani/pkg/reply"
	"coinbani/pkg/telegram"
	"coinbani/pkg/template"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg options.Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}

	logger, err := options.GetLogger(cfg.Log)
	defer logger.Sync()
	logger.Debug("initializing coinbani bot whit config", zap.String("cfg", fmt.Sprintf("%+v", cfg)))

	// Setup telegram bot
	bot, err := telegram.NewBot(cfg.Bot, logger)
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
