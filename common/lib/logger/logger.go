package logger

import (
	"push/common/lib/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	Raw *zap.Logger
}

func NewLogger(env env.Env) (*Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	switch env.Log.Level {
	case "debug":
		cfg.Level.SetLevel(zapcore.DebugLevel)
	case "info":
		cfg.Level.SetLevel(zapcore.InfoLevel)
	case "warn":
		cfg.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		cfg.Level.SetLevel(zapcore.ErrorLevel)
	case "fatal":
		cfg.Level.SetLevel(zapcore.FatalLevel)
	default:
		cfg.Level.SetLevel(zapcore.PanicLevel)
	}

	raw, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{
		SugaredLogger: raw.Sugar(),
		Raw:           raw,
	}, nil
}
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	newRaw := l.Raw.WithOptions(opts...)
	return &Logger{
		SugaredLogger: newRaw.Sugar(),
		Raw:           newRaw,
	}
}
