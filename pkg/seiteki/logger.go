package seiteki

// Logger describes a generic logger
// to be used for request logging
type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}

// loggerWrapper wraps around a given logger
// which checks if the wrapped logger is nil
// to simplify logger calls (so that logger != nil
// checks are not required).
type loggerWrapper struct {
	logger Logger
}

// newLoggerWrapper creates a new instance of
// loggerWrapper wrapping the given logger.
func newLogegrWrapper(logger Logger) *loggerWrapper {
	return &loggerWrapper{
		logger: logger,
	}
}

func (l *loggerWrapper) Error(args ...interface{}) {
	if l.logger != nil {
		l.logger.Error(args...)
	}
}

func (l *loggerWrapper) Errorf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Errorf(format, args...)
	}
}

func (l *loggerWrapper) Fatal(args ...interface{}) {
	if l.logger != nil {
		l.logger.Fatal(args...)
	}
}

func (l *loggerWrapper) Fatalf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalf(format, args...)
	}
}

func (l *loggerWrapper) Info(args ...interface{}) {
	if l.logger != nil {
		l.logger.Info(args...)
	}
}

func (l *loggerWrapper) Infof(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Infof(format, args...)
	}
}

func (l *loggerWrapper) Warning(args ...interface{}) {
	if l.logger != nil {
		l.logger.Warning(args...)
	}
}

func (l *loggerWrapper) Warningf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Warningf(format, args...)
	}
}
