package server

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/model"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type ReviewService interface {
	GetReviewCards(
		ctx context.Context,
		groupId model.GroupId,
		settings model.ReviewSettings,
	) ([]*model.Card, error)
	AddReviewResults(
		ctx context.Context, groupId model.GroupId,
		answers []model.ReviewCardResult,
	) error
	GetCardsMarks(ctx context.Context, groupId model.GroupId) ([]model.CardMark, error)
}

type ReviewServerDeps struct {
	ReviewService ReviewService
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
	groupId := model.GroupId(req.GroupId)

	settings := model.ReviewSettings{
		CardsCount: int(req.CardsCount),
	}

	cards, err := r.ReviewService.GetReviewCards(ctx, groupId, settings)
	if err != nil {
		return nil, err
	}

	cardsResp := make([]*pb.Card, len(cards))
	for i, card := range cards {
		cardsResp[i] = cardToResponse(card)
	}

	return &pb.GetReviewCardsResponse{FullCards: cardsResp}, nil
}
func (r *ReviewServer) AddReviewResults(ctx context.Context, req *pb.AddReviewResultsRequest) (
	*emptypb.Empty,
	error,
) {
	groupId := model.GroupId(req.GroupId)

	answers := make([]model.ReviewCardResult, len(req.GetCardResults()))
	for i, result := range req.GetCardResults() {
		answers[i] = model.ReviewCardResult{
			Answer:   answerToModels(result.CardAnswer),
			CardId:   model.CardId(result.CardId),
			Duration: result.Duration.AsDuration(),
		}
	}

	err := r.ReviewService.AddReviewResults(ctx, groupId, answers)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (r *ReviewServer) GetCardsProgress(
	ctx context.Context,
	req *pb.GetCardsProgressRequest,
) (*pb.GetCardsProgressResponse, error) {
	groupId := model.GroupId(req.GroupId)
	cardsMarks, err := r.ReviewService.GetCardsMarks(ctx, groupId)
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
