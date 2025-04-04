package verification

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"strings"
)

const (
	userIdKey   = "user-id"
	authMetaKey = "authorization"
	authBearer  = "Bearer"
)

type Verifier interface {
	VerifyToken(ctx context.Context, token string) (int64, error)
}

func Interceptor(log *slog.Logger, verificationService Verifier) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		token, err := getTokenFormContext(ctx)
		if err != nil {
			return nil, err
		}

		userId, err := verificationService.VerifyToken(ctx, token)
		if err != nil {
			log.Error("failed to verify token", sl.Err(err))
			return nil, fmt.Errorf("verification err: %w", err)
		}

		ctx = setUserId(ctx, entity.UserId(userId))

		return handler(ctx, req)
	}
}

func GetUserId(ctx context.Context) (entity.UserId, error) {
	userId, ok := ctx.Value(userIdKey).(entity.UserId)
	if !ok {
		return entity.UserId(0), entity.ErrVerificationFailed
	}
	return userId, nil
}

func setUserId(ctx context.Context, userId entity.UserId) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func getTokenFormContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", entity.ErrMetadataIsEmpty
	}

	values := md[authMetaKey]
	if len(values) == 0 {
		return "", entity.ErrNoAuthHeader
	}

	parts := strings.SplitN(values[0], " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], authBearer) {
		return "", entity.ErrIncorrectAuthHeader
	}

	return parts[1], nil
}
