package grpc

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ReviewUseCase interface {
	GetReviewCards(ctx context.Context, id entity.GroupId) ([]entity.Card, error)
	GetGroupReviewInfo(ctx context.Context, id entity.GroupId) (entity.GroupReviewInfo, error)
	UpdateGroupReviewInfo(ctx context.Context, reviewInfo entity.UpdateGroupReviewInfo) error
	SaveResults(ctx context.Context, answers []entity.ReviewCardResult) error
}
type ReviewServiceDeps struct {
	ReviewUseCase ReviewUseCase
}

type ReviewService struct {
	ReviewServiceDeps
	pb.UnimplementedReviewServiceServer
}

func NewReviewService(deps ReviewServiceDeps) *ReviewService {
	return &ReviewService{ReviewServiceDeps: deps}
}

func (r *ReviewService) GetGroupReviewInfo(
	ctx context.Context,
	req *pb.GetGroupReviewInfoRequest,
) (
	*pb.GetGroupReviewInfoResponse,
	error,
) {
	groupId := entity.GroupId(req.GroupId)

	reviewInfo, err := r.ReviewUseCase.GetGroupReviewInfo(ctx, groupId)
	if err != nil {
		return nil, err
	}
	return &pb.GetGroupReviewInfoResponse{CardsCount: int32(reviewInfo.CardsCount)}, nil
}

func (r *ReviewService) UpdateGroupReviewInfo(ctx context.Context, req *pb.UpdateGroupReviewInfoRequest) (
	*emptypb.Empty,
	error,
) {
	groupId := entity.GroupId(req.GroupId)

	reviewInfo := entity.UpdateGroupReviewInfo{GroupId: groupId, CardsCount: int(req.CardsCount)}

	err := r.ReviewUseCase.UpdateGroupReviewInfo(ctx, reviewInfo)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (r *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewResponse, error) {
	groupId := entity.GroupId(req.GroupId)

	cards, err := r.ReviewUseCase.GetReviewCards(ctx, groupId)
	if err != nil {
		return nil, err
	}

	cardsResp := make([]*pb.Card, len(cards))
	for i, card := range cards {
		cardsResp[i] = cardToResponse(card)
	}

	return &pb.GetReviewResponse{Cards: cardsResp}, nil
}
func (r *ReviewService) SaveReview(ctx context.Context, req *pb.CompleteReviewRequest) (*emptypb.Empty, error) {
	answers := make([]entity.ReviewCardResult, len(req.GetCardResults()))
	for i, result := range req.GetCardResults() {
		answers[i] = entity.ReviewCardResult{
			Answer:   answerToEntity(result.CardAnswer),
			CardId:   entity.CardId(result.CardId),
			Duration: result.Duration.AsDuration(),
		}
	}

	err := r.ReviewUseCase.SaveResults(ctx, answers)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
