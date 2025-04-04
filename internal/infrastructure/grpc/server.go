package grpc

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type groupUseCase interface {
	Create(ctx context.Context, name, description string, visibility entity.GroupVisibility) (entity.GroupId, error)
	List(ctx context.Context) ([]entity.Group, error)
	Get(ctx context.Context, id entity.GroupId) (entity.Group, error)
	Update(ctx context.Context, updateGroup entity.UpdateGroup) error
	Delete(ctx context.Context, id entity.GroupId) error
}

type cardsUseCase interface {
	Create(ctx context.Context, groupId entity.GroupId, frontText, backText string) (entity.CardId, error)
	List(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error)
	Get(ctx context.Context, id entity.CardId) (entity.Card, error)
	Update(ctx context.Context, card entity.UpdateCard) error
	Delete(ctx context.Context, id entity.CardId) error
}

type Server struct {
	pb.UnimplementedLearningCardsServer
	*slog.Logger
	groupUseCase
	cardsUseCase
}

func NewServer(groupUseCase groupUseCase, cardsUseCase cardsUseCase, logger *slog.Logger) *Server {
	return &Server{
		groupUseCase: groupUseCase,
		cardsUseCase: cardsUseCase,
		Logger:       logger,
	}
}

func (s *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (
	*pb.CreateGroupResponse,
	error,
) {

	visibility := entity.GroupVisibility(req.Visibility)
	groupId, err := s.groupUseCase.Create(ctx, req.GroupName, req.Description, visibility)

	if err != nil {
		return nil, err
	}

	resGroupId := int64(groupId)
	s.Logger.Info("Cards group created", slog.String("name", req.GroupName), slog.Int64("groupId", resGroupId))
	return &pb.CreateGroupResponse{GroupId: resGroupId}, nil
}

func (s *Server) ListGroups(ctx context.Context, _ *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	groups, err := s.groupUseCase.List(ctx)

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

	group, err := s.groupUseCase.Get(ctx, groupId)

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

	if err := s.groupUseCase.Update(ctx, group); err != nil {
		return nil, err
	}

	s.Logger.Info("Group updated", slog.Int64("groupId", req.GroupId))
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*emptypb.Empty, error) {
	groupId := entity.GroupId(req.GroupId)

	if err := s.groupUseCase.Delete(ctx, groupId); err != nil {
		return nil, err
	}

	s.Logger.Info("Group deleted", slog.Int64("groupId", req.GroupId))
	return &emptypb.Empty{}, nil
}

func (s *Server) AddCard(ctx context.Context, req *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	groupId := entity.GroupId(req.GroupId)
	cardId, err := s.cardsUseCase.Create(ctx, groupId, req.FrontText, req.BackText)

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
	cards, err := s.cardsUseCase.List(ctx, groupId)

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
	card, err := s.cardsUseCase.Get(ctx, cardId)

	if err != nil {
		return nil, err
	}

	cardResp := cardToResponse(card)
	return &pb.GetCardResponse{Card: cardResp}, nil
}

func (s *Server) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*emptypb.Empty, error) {
	cardId := entity.CardId(req.CardId)
	card := entity.UpdateCard{Id: cardId, FrontText: req.FrontText, BackText: req.BackText}

	if err := s.cardsUseCase.Update(ctx, card); err != nil {
		return nil, err
	}

	s.Logger.Info("Card updated", slog.Int64("cardId", req.CardId))
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*emptypb.Empty, error) {
	cardId := entity.CardId(req.CardId)

	if err := s.cardsUseCase.Delete(ctx, cardId); err != nil {
		return nil, err
	}

	s.Logger.Info("Card deleted", slog.Int64("cardId", req.CardId))
	return &emptypb.Empty{}, nil
}
