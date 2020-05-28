package main

import (
	"coinbani/pkg/template"
	"fmt"
	"log"
	"net/http"
	"time"

	"coinbani/cmd/coinbani/options"
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

	// Setup telegram bot
	bot.Debug = cfg.Bot.Debug
	logger.Info(fmt.Sprintf("authorized on account %s", bot.Self.UserName))

	logger.Info("setting up webhook", zap.String("CallbackURL", cfg.Application.CallbackURL))
	_, err = bot.SetWebhook(tb.NewWebhook(cfg.Application.CallbackURL + bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		logger.Error(fmt.Sprintf("Telegram callback failed: %s", info.LastErrorMessage))
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(":"+cfg.Application.Port, nil)

	// setup services
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
