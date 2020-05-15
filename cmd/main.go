package main

import (
	"coinbani/cmd/options"
	"fmt"
	"log"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("initializing coinbani bot")

	cfg := options.NewConfig()
	bot, err := tb.NewBotAPI(cfg.BotToken)
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

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				logger.Info("channel closed, exiting")
				break
			}

			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

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

			bot.Send(msg)
		}
	}

}
