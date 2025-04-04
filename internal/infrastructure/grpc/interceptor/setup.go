package interceptor

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log/slog"
)

const requestIdKey = "requestId"
const unknownValue = "unknown"

func SetupInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		requestId := uuid.New().String()
		ctx = context.WithValue(ctx, requestIdKey, requestId)
		return handler(ctx, req)
	}
}

func getRequestIdAttr(ctx context.Context) slog.Attr {
	requestId, ok := ctx.Value(requestIdKey).(string)
	if !ok {
		requestId = unknownValue
	}
	return slog.String(requestIdKey, requestId)
}
