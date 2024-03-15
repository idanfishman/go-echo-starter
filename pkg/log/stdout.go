package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitializeLogger initializes a new zap logger with the provided log level.
// The logger is configured to write JSON structured logs to stdout/stderr,
// with timestamps in RFC3339Nano format.
// Log example: {"level":"info","timestamp":"2021-08-25T14:00:00.000000000Z","message":"Fishman!"}
func InitializeLogger(logLevel zap.AtomicLevel, disableCaller bool, disableStacktrace bool) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.MessageKey = "message"
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339Nano))
	}

	config := zap.Config{
		Level:             logLevel,
		Development:       false,
		DisableCaller:     disableCaller,
		DisableStacktrace: disableStacktrace,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	return zap.Must(config.Build())
}

// UpdateLogLevel changes the log level of the provided zap.AtomicLevel to the target level.
// If the target level string cannot be parsed, the function will silently fail.
func UpdateLogLevel(currentLevel zap.AtomicLevel, targetLevel string) {
	l, _ := zapcore.ParseLevel(targetLevel)
	currentLevel.SetLevel(l)
}

// FlushLogger writes any buffered log entries to the output.
// The function should be called before the program exits to ensure all logs are written.
func FlushLogger(logger *zap.Logger) {
	_ = logger.Sync()
}
