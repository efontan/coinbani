package options

import (
	"log"
	"os"
	"strconv"
)

type config struct {
	Bot       *BotConfig
	Providers *ProvidersConfig
	Template  *TemplateConfig
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
		Providers: &ProvidersConfig{
			BBURL:         os.Getenv("BB_URL"),
			SatoshiARSURL: os.Getenv("SATOSHI_ARS_URL"),
			SatoshiUSDURL: os.Getenv("SATOSHI_USD_URL"),
			DollarURL:     os.Getenv("DOLLAR_URL"),
		},
		Template: &TemplateConfig{
			TemplatesDir: os.Getenv("TEMPLATES_DIR"),
		},
	}
}

type BotConfig struct {
	Token string
	Debug bool
}

type ProvidersConfig struct {
	BBURL         string
	SatoshiARSURL string
	SatoshiUSDURL string
	DollarURL     string
}

type TemplateConfig struct {
	TemplatesDir string
}
