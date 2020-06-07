package options

import (
	"go.uber.org/zap"
)

func GetLogger(c *LogConfig) (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()
	if c.Level == "debug" {
		logConfig = zap.NewDevelopmentConfig()
	}
	logConfig.Encoding = "json"
	return logConfig.Build()
}
