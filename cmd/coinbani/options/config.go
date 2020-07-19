package options

import (
	"encoding/json"
)

type Config struct {
	Bot       *BotConfig
	Log       *LogConfig
	Providers *ProvidersConfig
}

type BotConfig struct {
	CallbackURL      string `env:"CALLBACK_URL"`
	Debug            bool   `env:"BOT_DEBUG,default=false"`
	Port             string `env:"PORT"`
	Token            string `env:"BOT_TOKEN,required"`
	IsWebhookEnabled bool   `env:"BOT_IS_WEBHOOK_ENABLED,default=false"`
}

type ProvidersConfig struct {
	BBURL           string  `env:"BB_URL"`
	DollarURL       string  `env:"DOLLAR_URL"`
	DollarSavingTax float64 `env:"DOLLAR_SAVING_TAX"`
	SatoshiARSURL   string  `env:"SATOSHI_ARS_URL"`
	SatoshiUSDURL   string  `env:"SATOSHI_USD_URL"`
}

type LogConfig struct {
	Level string `env:"LOG_LEVEL,default=info"`
}

func (c *Config) String() string {
	res, _ := json.Marshal(c)
	return string(res)
}

func NewConfig() *Config {
	// savingTax, err := strconv.ParseFloat(os.Getenv("DOLLAR_SAVING_TAX"), 64)
	// if err != nil {
	// 	log.Panic(err)
	// }
	return &Config{}
}
