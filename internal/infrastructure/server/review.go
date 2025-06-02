package server

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type ReviewService interface {
	GetReviewCards(
		ctx context.Context,
		userId entity.UserId,
		groupId entity.GroupId,
		settings entity.ReviewSettings,
	) ([]entity.Card, error)
	AddReviewResults(
		ctx context.Context, userId entity.UserId, groupId entity.GroupId,
		answers []entity.ReviewCardResult,
	) error
	GetCardsMarks(ctx context.Context, userId entity.UserId, groupId entity.GroupId) ([]entity.CardMark, error)
}

type ReviewServerDeps struct {
	ReviewService ReviewService
	AuthVerifier  AuthVerifier
	Logger        *slog.Logger
}

type ReviewServer struct {
	ReviewServerDeps
	pb.UnimplementedReviewServiceServer
}

func NewReviewServer(deps ReviewServerDeps) *ReviewServer {
	return &ReviewServer{ReviewServerDeps: deps}
}

func (r *ReviewServer) GetReviewCards(ctx context.Context, req *pb.GetReviewCardsRequest) (
	*pb.GetReviewCardsResponse,
	error,
) {
	userId, err := r.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groupId := entity.GroupId(req.GroupId)

	settings := entity.ReviewSettings{
		CardsCount: int(req.CardsCount),
	}

	cards, err := r.ReviewService.GetReviewCards(ctx, userId, groupId, settings)
	if err != nil {
		return nil, err
	}

	cardsResp := make([]*pb.ReviewCard, len(cards))
	for i, card := range cards {
		cardsResp[i] = cardToReviewResponse(card)
	}

	return &pb.GetReviewCardsResponse{Cards: cardsResp}, nil
}
func (r *ReviewServer) AddReviewResults(ctx context.Context, req *pb.AddReviewResultsRequest) (
	*emptypb.Empty,
	error,
) {
	userId, err := r.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groupId := entity.GroupId(req.GroupId)

	answers := make([]entity.ReviewCardResult, len(req.GetCardResults()))
	for i, result := range req.GetCardResults() {
		answers[i] = entity.ReviewCardResult{
			Answer:   answerToEntity(result.CardAnswer),
			CardId:   entity.CardId(result.CardId),
			Duration: result.Duration.AsDuration(),
		}
	}

	err = r.ReviewService.AddReviewResults(ctx, userId, groupId, answers)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (r *ReviewServer) GetCardsProgress(
	ctx context.Context,
	req *pb.GetCardsProgressRequest,
) (*pb.GetCardsProgressResponse, error) {
	userId, err := r.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	groupId := entity.GroupId(req.GroupId)

	cardsMarks, err := r.ReviewService.GetCardsMarks(ctx, userId, groupId)
	if err != nil {
		return nil, err
	}
	cardsResp := make([]*pb.CardProgress, len(cardsMarks))
	for i, cardMark := range cardsMarks {
		cardsResp[i] = &pb.CardProgress{
			CardId: int64(cardMark.Id),
			Mark:   markToResponse(cardMark.Mark),
		}
	}
	return &pb.GetCardsProgressResponse{Cards: cardsResp}, nil
}
