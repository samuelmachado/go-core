package log

import (
	"context"

	"emperror.dev/errors"
	colorable "github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap wraps a zap.Logger and implements Logger inteface.
type Zap struct {
	logger *zap.Logger
}

var _ Logger = (*Zap)(nil)

// ZapConfig handle the config information that will be passed to zap.
type ZapConfig struct {
	Version           string
	DisableStackTrace bool
	Debug             bool `env:"DEBUG"`
}

func NewLoggerZap(config ZapConfig) (*Zap, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	loggerConfig.DisableStacktrace = config.DisableStackTrace
	loggerConfig.InitialFields = map[string]interface{}{
		"version": config.Version,
	}
	if config.Debug {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := loggerConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, errors.Wrap(err, "error on building zap logger")
	}

	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger = zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(loggerConfig.EncoderConfig),
			zapcore.AddSync(colorable.NewColorableStdout()),
			loggerConfig.Level,
		),
	)

	return &Zap{
		logger: logger,
	}, nil
}

func fieldsToZap(ctx context.Context, fs []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fs), len(fs)+1)

	for i := range fs {
		zapFields[i] = zap.Any(fs[i].Key, fs[i].Value)
	}

	return zapFields
}

// Debug will write a log with level debug.
func (z *Zap) Debug(ctx context.Context, msg string, fields ...Field) {
	z.logger.Debug(msg, fieldsToZap(ctx, fields)...)
}

// Error will write a log with level error.
func (z *Zap) Error(ctx context.Context, msg string, fields ...Field) {
	z.logger.Error(msg, fieldsToZap(ctx, fields)...)
}

// Fatal will write a log with level fatal.
func (z *Zap) Fatal(ctx context.Context, msg string, fields ...Field) {
	z.logger.Fatal(msg, fieldsToZap(ctx, fields)...)
}

// Info will write a log with level info.
func (z *Zap) Info(ctx context.Context, msg string, fields ...Field) {
	z.logger.Info(msg, fieldsToZap(ctx, fields)...)
}

// Panic will write a log with level panic.
func (z *Zap) Panic(ctx context.Context, msg string, fields ...Field) {
	z.logger.Panic(msg, fieldsToZap(ctx, fields)...)
}

// Warn will write a log with level warn.
func (z *Zap) Warn(ctx context.Context, msg string, fields ...Field) {
	z.logger.Warn(msg, fieldsToZap(ctx, fields)...)
}
