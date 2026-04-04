package loggers

import (
	"go.uber.org/zap"
)

type Loggers struct {
	HttpLogger *zap.Logger
}

func NewLoggers() *Loggers {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	baseLogger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return &Loggers{
		HttpLogger: baseLogger.With(
			zap.String("layer", "http"),
		),
	}
}
