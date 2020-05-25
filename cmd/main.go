package main

import (
	"coinbani/pkg/template"
	"fmt"
	"log"
	"net/http"
	"time"

	"coinbani/cmd/options"
	"coinbani/pkg/cache"
	"coinbani/pkg/currency"
	"coinbani/pkg/currency/provider"
	"coinbani/pkg/reply"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("initializing coinbani bot")

	cfg := options.NewConfig()
	bot, err := tb.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = cfg.Bot.Debug
	logger.Info(fmt.Sprintf("authorized on account %s", bot.Self.UserName))

	u := tb.NewUpdate(0)
	u.Timeout = 60

	logger.Info("starting channel for getting bot updates")
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	providerCache := cache.New()
	httpClient := &http.Client{Timeout: 10 * time.Second}
	bbProvider := provider.NewBBProvider(cfg.Providers, httpClient, providerCache)
	satoshiTProvider := provider.NewSatoshiTProvider(cfg.Providers, httpClient, providerCache)
	dollarProvider := provider.NewDollarProvider(cfg.Providers, httpClient, providerCache)

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
