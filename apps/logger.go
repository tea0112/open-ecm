package apps

import (
	"context"
	"log"

	"go.uber.org/zap"
)

type loggerKey struct{}

func CtxWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func LoggerFromCtx(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		panic(ctx)
	}
	return logger
}

func NewLogger(appMode string) *zap.Logger {
	var logger *zap.Logger
	var err error
	if appMode == "prod" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
