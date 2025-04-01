package interceptor

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

const requestIdKey = "requestId"
const userIdKey = "userId"
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

		md, ok := metadata.FromIncomingContext(ctx)
		userId := "unknown"
		if ok {
			if id := md.Get("user-id"); len(id) > 0 {
				userId = id[0]
			}
		}
		ctx = context.WithValue(ctx, userIdKey, userId)

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

func getUserIdAttr(ctx context.Context) slog.Attr {
	userId, ok := ctx.Value(userIdKey).(string)
	if !ok {
		userId = unknownValue
	}
	return slog.String(userIdKey, userId)
}
