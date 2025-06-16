package logger

import (
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type FxLogger struct {
	*Logger
}

func NewFxLogger(lc *Logger) fxevent.Logger {
	return &FxLogger{Logger: lc.WithOptions(zap.WithCaller(false))}
}

func (l *FxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Debugf("OnStart hook executing: callee=%s caller=%s", e.FunctionName, e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Errorf("OnStart hook failed: callee=%s caller=%s err=%v", e.FunctionName, e.CallerName, e.Err)
		} else {
			l.Debugf("OnStart hook executed: callee=%s caller=%s runtime=%v", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.OnStopExecuting:
		l.Debugf("OnStop hook executing: callee=%s caller=%s", e.FunctionName, e.CallerName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Errorf("OnStop hook failed: callee=%s caller=%s err=%v", e.FunctionName, e.CallerName, e.Err)
		} else {
			l.Debugf("OnStop hook executed: callee=%s caller=%s runtime=%v", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.Supplied:
		l.Debugf("supplied: type=%s err=%v", e.TypeName, e.Err)
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Debugf("provided: %s => %s", e.ConstructorName, rtype)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Infof("decorated: %s => %s", e.DecoratorName, rtype)
		}
	case *fxevent.Invoking:
		l.Debugf("invoking: %s", e.FunctionName)
	case *fxevent.Started:
		if e.Err == nil {
			l.Info("fx started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err == nil {
			l.Infof("custom fxevent.Logger initialized: %s", e.ConstructorName)
		}
	}
}
