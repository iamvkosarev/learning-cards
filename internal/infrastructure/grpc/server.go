package grpc

import (
	"context"
	"errors"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type GroupUseCase interface {
	Create(ctx context.Context, name string) (int64, error)
}

type Server struct {
	pb.UnimplementedLearningCardsServer
	*slog.Logger
	groupUseCase GroupUseCase
}

func NewServer(groupUseCase GroupUseCase, logger *slog.Logger) *Server {
	return &Server{groupUseCase: groupUseCase, Logger: logger}
}

func (s *Server) CreateCardsGroup(ctx context.Context, req *pb.CreateCardsGroupRequest) (
	*pb.CreateCardsGroupResponse,
	error,
) {
	const op = "grpc.Server.CreateCardsGroup"

	log := s.Logger.With(
		slog.String("op", op),
	)

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.groupUseCase.Create(ctx, req.GroupName)

	if err != nil {
		var verificationErr entity.VerificationError
		switch {
		case errors.Is(err, entity.ErrMetadataIsEmpty):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.As(err, &verificationErr):
			if verificationErr.StatusCode != codes.PermissionDenied &&
				verificationErr.StatusCode != codes.InvalidArgument {
				s.Logger.Error("failed to verify user", sl.Err(verificationErr))
				return nil, status.Error(codes.Internal, "internal verification error")
			} else {
				return nil, status.Error(verificationErr.StatusCode, verificationErr.Message)
			}

		}
		s.Logger.Info("failed to create group", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	log.Info("Cards group created", slog.String("name", req.GroupName), slog.Int64("id", id))
	return &pb.CreateCardsGroupResponse{GroupId: id}, nil
}
