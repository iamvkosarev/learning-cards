package interceptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor/verification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

const internalErrMessage = "internal server error"
const verificationErrMessage = "verification error"
const validationErrMessage = "not correct input"
const metadataEmptyErrMessage = "metadata is empty"
const cardNotFoundMessage = "card not found"
const groupNotFoundMessage = "group not found"

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
			err = mapDomainErrorToGRPC(err)
		}
		return resp, err
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
	return string(userId)
}

func mapDomainErrorToGRPC(err error) error {
	var verificationErr *entity.VerificationError
	if errors.As(err, &verificationErr) {
		return status.Error(codes.Internal, verificationErrMessage)
	}

	var validationErr *entity.ValidationError
	if errors.As(err, &validationErr) {
		return status.Error(codes.InvalidArgument, validationErrMessage)
	}

	switch {
	case errors.Is(err, entity.ErrMetadataIsEmpty):
		return status.Error(codes.InvalidArgument, metadataEmptyErrMessage)
	case errors.Is(err, entity.ErrNoAuthHeader):
		return status.Error(codes.InvalidArgument, verificationErrMessage)
	case errors.Is(err, entity.ErrIncorrectAuthHeader):
		return status.Error(codes.Unauthenticated, verificationErrMessage)
	case errors.Is(err, entity.ErrVerificationFailed):
		return status.Error(codes.Unauthenticated, verificationErrMessage)

	case errors.Is(err, entity.ErrGroupNotFound):
		return status.Error(codes.NotFound, cardNotFoundMessage)
	case errors.Is(err, entity.ErrCardNotFound):
		return status.Error(codes.NotFound, groupNotFoundMessage)
	default:
		return status.Error(codes.Internal, internalErrMessage)
	}
}
