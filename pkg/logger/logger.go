package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger struct {
	internal *zap.Logger
}

func NewLogger() (*Logger, error) {
	conf := zap.NewDevelopmentConfig()
	conf.OutputPaths = []string{"stdout"}

	log, err := conf.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{
		internal: log,
	}, nil
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = append(contextFields(ctx), fields...)
	}

	l.internal.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = append(contextFields(ctx), fields...)
	}

	l.internal.Warn(msg, fields...)
}

func (l *Logger) WithOptions(option ...zap.Option) *Logger {
	return &Logger{
		internal: l.internal.WithOptions(option...),
	}
}

func contextFields(ctx context.Context) []zap.Field {
	userID := ctx.Value("userId")
	if userID != nil {
		return []zap.Field{
			zap.Any("userId", userID),
		}
	}

	return []zap.Field{}
}
