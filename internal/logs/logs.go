package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var log *zap.Logger

func InitializeLogger() *zap.SugaredLogger {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger:  " + err.Error())
	}
	sugar := logger.Sugar()
	zap.ReplaceGlobals(logger)

	return sugar

}

func SyncLogger(logger *zap.SugaredLogger) {
	_ = logger.Sync()
}
