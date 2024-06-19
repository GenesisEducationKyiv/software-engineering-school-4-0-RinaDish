package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
    Info(args ...interface{})
    Infof(template string, args ...interface{})
    Error(args ...interface{})
    Errorf(template string, args ...interface{})
    With(args ...interface{}) Logger

}

type ZapLogger struct {
    logger *zap.SugaredLogger
}

func NewZapLogger(sugaredLogger *zap.SugaredLogger) ZapLogger {
    return ZapLogger{logger: sugaredLogger}
}

func (l ZapLogger) Info(args ...interface{}) {
    l.logger.Info(args...)
}

func (l ZapLogger) Infof(template string, args ...interface{}) {
    l.logger.Infof(template, args...)
}

func (l ZapLogger) Error(args ...interface{}) {
    l.logger.Error(args...)
}

func (l ZapLogger) Errorf(template string, args ...interface{}) {
    l.logger.Errorf(template, args...)
}

func (l ZapLogger) With(args ...interface{}) Logger {
    return ZapLogger{logger: l.logger.With(args...)}
}
