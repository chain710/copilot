package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var config = zap.NewProductionConfig()

func init() {
	setup()
}

func setup(options ...zap.Option) {
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	l, err := config.Build(options...)
	if err != nil {
		panic(err)
	}

	logger = l.Sugar()
}

func SetLogLevel(l zapcore.Level) {
	config.Level.SetLevel(l)
}

func Sync() error {
	return logger.Sync()
}

func With(args ...any) *zap.SugaredLogger {
	return logger.With(args...)
}

func Debugf(template string, args ...any) {
	logger.Debugf(template, args...)
}

func Infof(template string, arg ...any) {
	logger.Infof(template, arg...)
}

func Warnf(template string, arg ...any) {
	logger.Warnf(template, arg...)
}

func Errorf(template string, arg ...any) {
	logger.Errorf(template, arg...)
}

func Fatalf(template string, arg ...any) {
	logger.Fatalf(template, arg...)
}

func Debugw(msg string, keysAndValues ...any) {
	logger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...any) {
	logger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...any) {
	logger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...any) {
	logger.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...any) {
	logger.Fatalw(msg, keysAndValues...)
}
