package logger

type GinLogger struct {
	*Logger
}

func NewGinLogger(base *Logger) *GinLogger {
	return &GinLogger{Logger: base}
}

// Implements gin's Writer interface
func (l *GinLogger) Write(p []byte) (n int, err error) {
	l.Info(string(p))
	return len(p), nil
}
