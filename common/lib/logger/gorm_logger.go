package logger

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	*Logger
	Config logger.Config
}

func NewGormLogger(base *Logger) *GormLogger {
	return &GormLogger{
		Logger: base,
		Config: logger.Config{
			SlowThreshold: time.Second * 3, // 3초 이상이면 느린 쿼리로 간주
			LogLevel:      logger.Info,
		},
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.Config.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Info {
		l.Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Warn {
		l.Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Error {
		l.Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Config.LogLevel == logger.Silent {
		return
	}

	sql, rows := fc()
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.Config.LogLevel >= logger.Error:
		l.Errorf("[%.2fms][rows:%v] %s | err: %v", float64(elapsed.Milliseconds()), rows, sql, err)
	case elapsed > l.Config.SlowThreshold && l.Config.LogLevel >= logger.Warn:
		l.Warnf("[SLOW %.2fms][rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)
	case l.Config.LogLevel >= logger.Info:
		l.Debugf("[%.2fms][rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)
	}
}
