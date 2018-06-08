package logging

import (
	"context"

	"os"
	"path"

	"go.uber.org/zap"
)

type correlationIdType int
const (
	requestIdKey correlationIdType = iota
)

var logger *zap.Logger

func init() {
	zlog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = zlog.With(
		zap.Int("pid", os.Getpid()),
		zap.String("exe", path.Base(os.Args[0])),
	)
}

func WithRqId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}

func Logger(ctx context.Context) *zap.Logger {
	newLogger := logger

	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("requestId", ctxRqId))
		}
	}

	return newLogger
}