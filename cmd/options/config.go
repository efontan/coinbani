package options

import "os"

type config struct {
	BotToken string
}

func NewConfig() *config {
	return &config{
		BotToken: os.Getenv("BOT_TOKEN"),
	}
}
