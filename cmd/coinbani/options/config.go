package options

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Application *ApplicationConfig
	Bot         *BotConfig
	Providers   *ProvidersConfig
	Log         *LogConfig
}

func NewConfig() *Config {
	debug, err := strconv.ParseBool(os.Getenv("BOT_DEBUG"))
	if err != nil {
		log.Panic(err)
	}

	savingTax, err := strconv.ParseFloat(os.Getenv("DOLLAR_SAVING_TAX"), 64)
	if err != nil {
		log.Panic(err)
	}

	return &Config{
		Application: &ApplicationConfig{
			CallbackURL: os.Getenv("CALLBACK_URL"),
			Port:        os.Getenv("PORT"),
		},
		Bot: &BotConfig{
			Token: os.Getenv("BOT_TOKEN"),
			Debug: debug,
		},
		Providers: &ProvidersConfig{
			BBURL:           os.Getenv("BB_URL"),
			SatoshiARSURL:   os.Getenv("SATOSHI_ARS_URL"),
			SatoshiUSDURL:   os.Getenv("SATOSHI_USD_URL"),
			DollarURL:       os.Getenv("DOLLAR_URL"),
			DollarSavingTax: savingTax,
		},
		Log: &LogConfig{
			Level: os.Getenv("LOG_LEVEL"),
		},
	}
}

type ApplicationConfig struct {
	CallbackURL string
	Port        string
}

type BotConfig struct {
	Token string
	Debug bool
}

type ProvidersConfig struct {
	BBURL           string
	SatoshiARSURL   string
	SatoshiUSDURL   string
	DollarURL       string
	DollarSavingTax float64
}

type LogConfig struct {
	Level string
}

func (c *Config) String() string {
	res, _ := json.Marshal(c)
	return string(res)
}
