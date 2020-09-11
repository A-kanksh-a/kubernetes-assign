package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/controlplaneio/kubesec/pkg/ruler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger
)

func main() {

	var err error

	// logger writes to stderr
	logger, err = NewLogger("info", "console")
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	body, err := ioutil.ReadFile("deployment.yaml")
	if err != nil {

	}
	reports, err := ruler.NewRuleset(logger).Run(body)
	for _, report := range reports {
		// FIXME: encoding as json now for integrity check, this is the wrong way
		// to compute the hash over the result. Also, some error checking would be
		// more than ideal.
		reportValue, _ := json.Marshal(report)
		f, err := os.OpenFile("abc.json", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.Write(reportValue); err != nil {
			panic(err)
		}

	}
}

func NewLogger(logLevel string, zapEncoding string) (*zap.SugaredLogger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         zapEncoding,
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
