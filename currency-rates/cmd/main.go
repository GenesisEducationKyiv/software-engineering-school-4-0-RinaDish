package main

import (
	"context"
	"log"
	"time"

	"github.com/RinaDish/currency-rates/internal/app"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var cfg app.Config
	err := envconfig.Process("currency-rates", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	l, _ := config.Build()
	defer l.Sync()
	logger := l.Sugar()

	ctx := context.Background()

	application := app.NewApp(cfg, logger)
	if err := application.Run(ctx); err != nil {
		logger.Error(err)
	}
}
