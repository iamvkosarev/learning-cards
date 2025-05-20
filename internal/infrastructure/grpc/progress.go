package grpc

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
)

type ProgressUseCase interface {
	ListGroupsProgress(ctx context.Context) ([]entity.GroupProgress, error)
	GetCardsProgress(ctx context.Context, id entity.GroupId) ([]entity.CardProgress, error)
}

type ProgressServiceDeps struct {
	ProgressUseCase ProgressUseCase
}

type ProgressService struct {
	pb.UnimplementedProgressServiceServer
	ProgressServiceDeps
}

func NewProgressService(deps ProgressServiceDeps) *ProgressService {
	return &ProgressService{ProgressServiceDeps: deps}
}

func (p *ProgressService) ListGroupsProgress(
	ctx context.Context, _ *pb.ListGroupsProgressRequest,
) (*pb.ListGroupsProgressResponse, error) {
	groupsProgress, err := p.ProgressUseCase.ListGroupsProgress(ctx)
	if err != nil {
		return nil, err
	}
	groupProgressResp := make([]*pb.GroupProgress, len(groupsProgress))
	for i, prg := range groupsProgress {
		groupProgressResp[i] = groupProgressToResponse(prg)
	}

	return &pb.ListGroupsProgressResponse{GroupProgress: groupProgressResp}, nil
}

func (p *ProgressService) ListCardsProgress(
	ctx context.Context, req *pb.ListCardsProgressRequest,
) (
	*pb.ListCardsProgressResponse,
	error,
) {

	groupId := entity.GroupId(req.GroupId)

	cardsProgress, err := p.ProgressUseCase.GetCardsProgress(ctx, groupId)
	if err != nil {
		return nil, err
	}
	cardsProgressResp := make([]*pb.CardProgress, len(cardsProgress))
	for i, cardProgress := range cardsProgress {
		cardsProgressResp[i] = cardProgressToResponse(cardProgress)
	}

	return &pb.ListCardsProgressResponse{CardsProgress: cardsProgressResp}, nil
}
