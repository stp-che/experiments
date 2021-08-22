package behaviour

import "go.uber.org/zap"

var logger *zap.SugaredLogger = nil

func SetLogger(log *zap.SugaredLogger) {
	logger = log
}

func debugf(str string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Debugf(str, args...)
}

func fatalf(str string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Fatalf(str, args...)
}
