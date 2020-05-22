package options

import (
	"log"
	"os"
	"strconv"
)

type config struct {
	Bot      *BotConfig
	Exchange *ExchangeConfig
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
		Exchange: &ExchangeConfig{
			BBURL: os.Getenv("BB_URL"),
			SatoshiURL: os.Getenv("SATOSHI_URL"),
		},
	}
}

type BotConfig struct {
	Token string
	Debug bool
}

type ExchangeConfig struct {
	BBURL string
	SatoshiURL string
}
