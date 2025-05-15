package grpc

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type GroupUseCase interface {
	Create(ctx context.Context, name, description string, visibility entity.GroupVisibility) (entity.GroupId, error)
	List(ctx context.Context) ([]entity.Group, error)
	Get(ctx context.Context, id entity.GroupId) (entity.Group, error)
	Update(ctx context.Context, updateGroup entity.UpdateGroup) error
	Delete(ctx context.Context, id entity.GroupId) error
}

type CardsUseCase interface {
	Create(ctx context.Context, groupId entity.GroupId, frontText, backText string) (entity.CardId, error)
	List(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error)
	Get(ctx context.Context, id entity.CardId) (entity.Card, error)
	Update(ctx context.Context, card entity.UpdateCard) error
	Delete(ctx context.Context, id entity.CardId) error
}

type ProgressUseCase interface {
	ListGroupsProgress(ctx context.Context) ([]entity.GroupProgress, error)
	GetCardsProgress(ctx context.Context, id entity.GroupId) ([]entity.CardProgress, error)
	UpdateProgress(ctx context.Context, card []entity.ReviewCardResult) error
}

type ReviewUseCase interface {
	GetReviewCards(ctx context.Context, id entity.GroupId) ([]entity.Card, error)
	GetGroupReviewInfo(ctx context.Context, id entity.GroupId) (entity.GroupReviewInfo, error)
	UpdateGroupReviewInfo(ctx context.Context, reviewInfo entity.UpdateGroupReviewInfo) error
}

type Deps struct {
	GroupUseCase    GroupUseCase
	CardsUseCase    CardsUseCase
	ProgressUseCase ProgressUseCase
	ReviewUseCase   ReviewUseCase
	Logger          *slog.Logger
}

type Server struct {
	pb.UnimplementedLearningCardsServer
	Deps
}

func NewServer(deps Deps) *Server {
	return &Server{
		Deps: deps,
	}
}

func (s *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (
	*pb.CreateGroupResponse,
	error,
) {

	visibility := entity.GroupVisibility(req.Visibility)
	groupId, err := s.GroupUseCase.Create(ctx, req.GroupName, req.Description, visibility)

	if err != nil {
		return nil, err
	}

	resGroupId := int64(groupId)
	s.Logger.Info("Cards group created", slog.String("name", req.GroupName), slog.Int64("groupId", resGroupId))
	return &pb.CreateGroupResponse{GroupId: resGroupId}, nil
}

func (s *Server) ListGroups(ctx context.Context, _ *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	groups, err := s.GroupUseCase.List(ctx)

	if err != nil {
		return nil, err
	}

	var respGroups []*pb.CardsGroup
	for _, group := range groups {
		respGroups = append(
			respGroups, groupToResponse(group),
		)
	}

	return &pb.ListGroupsResponse{Groups: respGroups}, nil
}

func (s *Server) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	groupId := entity.GroupId(req.GroupId)

	group, err := s.GroupUseCase.Get(ctx, groupId)

	if err != nil {
		return nil, err
	}

	groupResp := groupToResponse(group)
	return &pb.GetGroupResponse{Group: groupResp}, nil
}

func (s *Server) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*emptypb.Empty, error) {
	groupId := entity.GroupId(req.GroupId)
	visibility := entity.GroupVisibility(req.Visibility)
	group := entity.UpdateGroup{
		Id:          groupId,
		Name:        req.GroupName,
		Description: req.Description,
		Visibility:  visibility,
	}

	if err := s.GroupUseCase.Update(ctx, group); err != nil {
		return nil, err
	}

	s.Logger.Info("Group updated", slog.Int64("groupId", req.GroupId))
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*emptypb.Empty, error) {
	groupId := entity.GroupId(req.GroupId)

	if err := s.GroupUseCase.Delete(ctx, groupId); err != nil {
		return nil, err
	}

	s.Logger.Info("Group deleted", slog.Int64("groupId", req.GroupId))
	return &emptypb.Empty{}, nil
}

func (s *Server) AddCard(ctx context.Context, req *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	groupId := entity.GroupId(req.GroupId)
	cardId, err := s.CardsUseCase.Create(ctx, groupId, req.FrontText, req.BackText)

	if err != nil {
		return nil, err
	}

	resId := int64(cardId)
	s.Logger.Info(
		"Card added", slog.Int64("cardId", resId),
		slog.Int64("groupId", req.GroupId),
	)
	return &pb.AddCardResponse{CardId: resId}, nil
}

func (s *Server) ListCards(ctx context.Context, req *pb.ListCardsRequest) (
	*pb.ListCardsResponse,
	error,
) {

	groupId := entity.GroupId(req.GroupId)
	cards, err := s.CardsUseCase.List(ctx, groupId)

	if err != nil {
		return nil, err
	}

	var respCards []*pb.Card
	for _, card := range cards {
		respCards = append(
			respCards, cardToResponse(card),
		)
	}

	return &pb.ListCardsResponse{Cards: respCards}, nil
}

func (s *Server) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	cardId := entity.CardId(req.CardId)
	card, err := s.CardsUseCase.Get(ctx, cardId)

	if err != nil {
		return nil, err
	}

	cardResp := cardToResponse(card)
	return &pb.GetCardResponse{Card: cardResp}, nil
}

func (s *Server) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*emptypb.Empty, error) {
	cardId := entity.CardId(req.CardId)
	card := entity.UpdateCard{Id: cardId, FrontText: req.FrontText, BackText: req.BackText}

	if err := s.CardsUseCase.Update(ctx, card); err != nil {
		return nil, err
	}

	s.Logger.Info("Card updated", slog.Int64("cardId", req.CardId))
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*emptypb.Empty, error) {
	cardId := entity.CardId(req.CardId)

	if err := s.CardsUseCase.Delete(ctx, cardId); err != nil {
		return nil, err
	}

	s.Logger.Info("Card deleted", slog.Int64("cardId", req.CardId))
	return &emptypb.Empty{}, nil
}

func (s *Server) ListGroupsProgress(ctx context.Context, req *pb.ListGroupsProgressRequest) (
	*pb.ListGroupsProgressResponse,
	error,
) {
	progresses, err := s.ProgressUseCase.ListGroupsProgress(ctx)
	if err != nil {
		return nil, err
	}
	progressesResp := make([]*pb.GroupProgress, len(progresses))
	for i, prg := range progresses {
		progressesResp[i] = groupProgressToResponse(prg)
	}

	return &pb.ListGroupsProgressResponse{GroupProgresses: progressesResp}, nil

}
func (s *Server) ListCardsProgress(ctx context.Context, req *pb.ListCardsProgressRequest) (
	*pb.ListCardsProgressResponse,
	error,
) {
	groupId := entity.GroupId(req.GroupId)

	cardsProgress, err := s.ProgressUseCase.GetCardsProgress(ctx, groupId)
	if err != nil {
		return nil, err
	}
	cardsProgressResp := make([]*pb.CardProgress, len(cardsProgress))
	for i, cardProgress := range cardsProgress {
		cardsProgressResp[i] = cardProgressToResponse(cardProgress)
	}

	return &pb.ListCardsProgressResponse{CardsProgresses: cardsProgressResp}, nil
}

func (s *Server) GetGroupReviewInfo(ctx context.Context, req *pb.GetGroupReviewInfoRequest) (
	*pb.GetGroupReviewInfoResponse,
	error,
) {
	groupId := entity.GroupId(req.GroupId)

	reviewInfo, err := s.ReviewUseCase.GetGroupReviewInfo(ctx, groupId)
	if err != nil {
		return nil, err
	}
	return &pb.GetGroupReviewInfoResponse{CardsCount: int32(reviewInfo.CardsCount)}, nil
}

func (s *Server) UpdateGroupReviewInfo(ctx context.Context, req *pb.UpdateGroupReviewInfoRequest) (
	*emptypb.Empty,
	error,
) {
	groupId := entity.GroupId(req.GroupId)

	reviewInfo := entity.UpdateGroupReviewInfo{GroupId: groupId, CardsCount: int(req.CardsCount)}

	err := s.ReviewUseCase.UpdateGroupReviewInfo(ctx, reviewInfo)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) MakeReview(ctx context.Context, req *pb.MakeReviewRequest) (*pb.MakeReviewResponse, error) {
	groupId := entity.GroupId(req.GroupId)

	cards, err := s.ReviewUseCase.GetReviewCards(ctx, groupId)
	if err != nil {
		return nil, err
	}

	cardsResp := make([]*pb.Card, len(cards))
	for i, card := range cards {
		cardsResp[i] = cardToResponse(card)
	}

	return &pb.MakeReviewResponse{Cards: cardsResp}, nil
}

func (s *Server) CompleteReview(ctx context.Context, req *pb.CompleteReviewRequest) (*emptypb.Empty, error) {
	answers := make([]entity.ReviewCardResult, len(req.GetCardResults()))
	for i, result := range req.GetCardResults() {
		answers[i] = entity.ReviewCardResult{
			Answer:   answerToEntity(result.CardAnswer),
			CardId:   entity.CardId(result.CardId),
			Duration: result.Duration.AsDuration(),
		}
	}

	err := s.ProgressUseCase.UpdateProgress(ctx, answers)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
