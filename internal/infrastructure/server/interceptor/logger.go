package interceptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

const internalErrMessage = "internal server error"
const userIdLogArg = "user_id"

func LoggerUnaryServerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		log := logger.With(
			getRequestIdAttr(ctx),
			getUserIdAttr(ctx),
			slog.String("method", info.FullMethod),
		)

		log.Info(
			"Processing request",
			slog.String("request_type", fmt.Sprintf("%T", req)),
		)

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		log = log.With(slog.Duration("duration", duration))

		if err != nil {
			st, _ := status.FromError(err)
			log.Error(
				"Request failed",
				slog.String("status", st.Code().String()),
				slog.String("response_message", st.Message()),
				sl.Err(err),
			)
		} else {
			log.Info(
				"Request succeeded",
				slog.String("status", "OK"),
				slog.String("response_type", fmt.Sprintf("%T", resp)),
			)
		}

		if err != nil {
			return resp, mapDomainErrorToGRPC(err)
		}
		return resp, nil
	}
}

func getUserIdAttr(ctx context.Context) slog.Attr {
	return slog.String(userIdLogArg, getUserId(ctx))
}

func getUserId(ctx context.Context) string {
	userId, err := verification.GetUserId(ctx)
	if err != nil {
		return unknownValue
	}

	return fmt.Sprintf("%d", userId)
}

func mapDomainErrorToGRPC(err error) error {
	var validationErr *model.ValidationError
	if errors.As(err, &validationErr) {
		return status.Error(codes.InvalidArgument, validationErr.Error())
	}

	var serverErr *model.ServerError
	if errors.As(err, &serverErr) {
		return status.Error(serverErr.Code, serverErr.Message)
	}
	return status.Error(codes.Internal, internalErrMessage)
}
