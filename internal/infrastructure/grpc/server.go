package grpc

import (
	"context"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type GroupUseCase interface {
	Create(ctx context.Context, name string, userId int64) (int64, error)
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

	userId := int64(0)

	id, err := s.groupUseCase.Create(ctx, req.GroupName, userId)

	log.Info("Cards group created", slog.String("name", req.GroupName), slog.Int64("id", id))
	return &pb.CreateCardsGroupResponse{GroupId: id}, nil
}
