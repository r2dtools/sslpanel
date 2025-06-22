package logger

import (
	"backend/config"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Error(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

type logger struct {
	zapLogger *zap.SugaredLogger
}

func (l *logger) Error(message string, args ...interface{}) {
	l.zapLogger.Errorf(message, args...)
}

func (l *logger) Warning(message string, args ...interface{}) {
	l.zapLogger.Warnf(message, args...)
}

func (l *logger) Info(message string, args ...interface{}) {
	l.zapLogger.Infof(message, args...)
}

func (l *logger) Debug(message string, args ...interface{}) {
	l.zapLogger.Debugf(message, args...)
}

func NewLogger(config *config.Config) (Logger, error) {
	var loggerConfig zap.Config
	outputPaths := []string{}

	if config.IsDevMode {
		loggerConfig = zap.NewDevelopmentConfig()
		outputPaths = append(outputPaths, "stderr")
	} else {
		err := createLogDirIfNeeded(config.LogFile)

		if err != nil {
			return nil, err
		}

		loggerConfig = zap.NewProductionConfig()
		outputPaths = append(outputPaths, config.LogFile)
	}

	loggerConfig.OutputPaths = outputPaths
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	zLogger, err := loggerConfig.Build()

	if err != nil {
		return nil, err
	}

	return &logger{zapLogger: zLogger.Sugar()}, nil
}

func createLogDirIfNeeded(logFilePath string) error {
	logDir := path.Dir(logFilePath)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}
