package main

import (
	"context"
	"log"
	"time"

	"github.com/RinaDish/currency-rates/internal/app"
	"github.com/RinaDish/currency-rates/tools"
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
	defer l.Sync() // nolint:errcheck
	sugaredLogger := l.Sugar()

	appLogger := tools.NewZapLogger(sugaredLogger)

	ctx := context.Background()

	application, err := app.NewApp(cfg, appLogger, ctx)
	if err != nil {
		appLogger.Error(err)
	}

	if err := application.Run(); err != nil {
		appLogger.Error(err)
	}
}
