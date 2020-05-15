package main

import (
	"fmt"
	"log"

	"coinbani/cmd/options"
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

	bot.Debug = true
	logger.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	u := tb.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	h := reply.NewHandler(bot, logger)

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				logger.Info("channel closed, exiting")
				break
			}

			go h.HandleReply(update)
		}
	}

}
