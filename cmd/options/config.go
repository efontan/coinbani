package options

import (
	"log"
	"os"
	"strconv"
)

type config struct {
	Bot *BotConfig
}

func NewConfig() *config {
	debug, err := strconv.ParseBool(os.Getenv("BOT_DEBUG"))
	if err != nil {
		log.Panic(err)
	}

	return &config{
		Bot: &BotConfig{
			Token: os.Getenv("BOT_TOKEN"),
			Debug: debug,
		},
	}
}

type BotConfig struct {
	Token string
	Debug bool
}
