package interceptor

import (
	"context"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc"
	"log/slog"
)

type Validator interface {
	Validate() error
}

func ValidationInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if validator, ok := req.(Validator); ok {
			if err := validator.Validate(); err != nil {

				logger.Warn(
					"Request validation failed",
					slog.String("method", info.FullMethod),
					getRequestIdAttr(ctx),
					getUserIdAttr(ctx),
					sl.Err(err),
				)
				return nil, entity.NewValidationError(err)
			}
		}

		return handler(ctx, req)
	}
}
