package logging

import (
	"context"
	"errors"

	"os"
	"path"

	"go.uber.org/zap"
	"github.com/google/uuid"
)

type correlationIdType int
const (
	requestIdKey correlationIdType = iota
	handlerIdKey
	methodIdKey
	endpointIdKey
	funcIdKey
	packageIdKey
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

func WithRqId(ctx context.Context, requestId string, funcHandler string, method string, url string) context.Context {
	c := context.WithValue(ctx, requestIdKey, requestId)

	if funcHandler != "" {
		c = context.WithValue( c, handlerIdKey, funcHandler )
	}

	if method != "" {
		c = context.WithValue( c, methodIdKey, method )
	}

	if url != "" {
		c = context.WithValue( c, endpointIdKey, url )
	}

	return c
}

func WithFuncId(ctx context.Context, funcName string, pkgName string ) context.Context {
	c := context.WithValue(ctx, funcIdKey, funcName)

	if pkgName != "" {
		c = context.WithValue( c, packageIdKey, pkgName )
	}

	return c
}

func Logger(ctx context.Context) *zap.Logger {
	newLogger := logger

	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("requestId", ctxRqId))
		}

		if handler, ok := ctx.Value(handlerIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("handler", handler))
		}

		if method, ok := ctx.Value(methodIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("method", method))
		}

		if url, ok := ctx.Value(endpointIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("endpoint", url))
		}

		if funcName, ok := ctx.Value(funcIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("function", funcName))
		}

		if pkgName, ok := ctx.Value(packageIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("package", pkgName))
		}
	}

	return newLogger
}

func GetReqId( ctx context.Context ) ( string, error ) {
	if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
		return ctxRqId, nil
	}

	rqId, err := uuid.NewRandom()
	if err != nil {
		return "", errors.New("requestId not found")
	}

	return rqId.String(), errors.New("requestId not found")
}