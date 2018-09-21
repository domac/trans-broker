package logger

import (
	"fmt"
	"go.uber.org/zap"
)

const (
	API_LOG_TIME_FORMAT = "2006-01-02 15:04:05.000000"
)

var defaultLogger *MyLogger

type MyLogger struct {
	zaplog *zap.Logger
}

func InitLogger(logFile string) (*MyLogger, error) {
	l := &MyLogger{}

	var err error

	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true

	if logFile != "" {
		cfg.OutputPaths = []string{logFile}
	}

	fmt.Printf("log path : %s\n", cfg.OutputPaths)

	l.zaplog, err = cfg.Build()

	if err != nil {
		return nil, err
	}

	defaultLogger = l

	return defaultLogger, nil
}

func GetLogger() *MyLogger {
	return defaultLogger
}

func GetInnerLogger() *zap.Logger {
	return defaultLogger.zaplog
}

func (l *MyLogger) Flush() error {
	return l.zaplog.Sync()
}

func (l *MyLogger) GetSugar() *zap.SugaredLogger {
	return l.zaplog.Sugar()
}

func Log() *zap.SugaredLogger {
	dlogger := GetLogger()
	if dlogger == nil {
		return nil
	}
	return dlogger.zaplog.Sugar()
}
