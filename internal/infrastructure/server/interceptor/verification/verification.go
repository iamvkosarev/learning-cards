package verification

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/model"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
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
	Close()
}

func Interceptor(log *slog.Logger, verificationService Verifier) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == pb.CardService_HealthCheck_FullMethodName {
			return nil, nil
		}
		token, err := GetTokenFormContext(ctx)
		if err != nil {
			return nil, err
		}

		userId, err := verificationService.VerifyToken(ctx, token)
		if err != nil {
			log.Error("failed to verify token", sl.Err(err))
			return nil, fmt.Errorf("verification err: %w", err)
		}

		ctx = setUserId(ctx, model.UserId(userId))

		return handler(ctx, req)
	}
}

func GetUserId(ctx context.Context) (model.UserId, error) {
	userId, ok := ctx.Value(userIdKey).(model.UserId)
	if !ok {
		return model.UserId(0), model.ErrVerificationFailed
	}
	return userId, nil
}

func setUserId(ctx context.Context, userId model.UserId) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func GetTokenFormContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", model.ErrMetadataIsEmpty
	}
	authValues := md[authMetaKey]
	if len(authValues) == 0 {
		return "", model.ErrNoAuthHeader
	}

	parts := strings.SplitN(authValues[0], " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], authBearer) {
		return "", model.ErrIncorrectAuthHeader
	}

	return parts[1], nil
}
