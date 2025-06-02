package server

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type GroupService interface {
	CreateGroup(
		ctx context.Context,
		userId entity.UserId,
		name, description string,
		visibility entity.GroupVisibility,
	) (
		entity.GroupId,
		error,
	)
	List(ctx context.Context, userId entity.UserId) ([]entity.Group, error)
	GetGroup(ctx context.Context, userId entity.UserId, id entity.GroupId) (entity.Group, error)
	UpdateGroup(ctx context.Context, userId entity.UserId, updateGroup entity.UpdateGroup) error
	DeleteGroup(ctx context.Context, userId entity.UserId, id entity.GroupId) error
}

type CardsService interface {
	Create(
		ctx context.Context, userId entity.UserId, groupId entity.GroupId, frontText,
		backText string,
	) (entity.CardId, error)
	ListCards(ctx context.Context, userId entity.UserId, groupId entity.GroupId) ([]entity.Card, error)
	GetCard(ctx context.Context, userId entity.UserId, id entity.CardId) (entity.Card, error)
	UpdateCard(ctx context.Context, userId entity.UserId, card entity.UpdateCard) error
	DeleteCard(ctx context.Context, userId entity.UserId, id entity.CardId) error
}

type CardServerDeps struct {
	GroupService GroupService
	CardsService CardsService
	AuthVerifier AuthVerifier
	Logger       *slog.Logger
}

type CardServer struct {
	pb.UnimplementedCardServiceServer
	CardServerDeps
}

func NewCardServer(deps CardServerDeps) *CardServer {
	return &CardServer{
		CardServerDeps: deps,
	}
}

func (s *CardServer) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (
	*pb.CreateGroupResponse,
	error,
) {
	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	visibility := entity.GroupVisibility(req.Visibility)
	groupId, err := s.GroupService.CreateGroup(ctx, userId, req.GroupName, req.Description, visibility)

	if err != nil {
		return nil, err
	}

	resGroupId := int64(groupId)
	s.Logger.Info("Cards group created", slog.String("name", req.GroupName), slog.Int64("groupId", resGroupId))
	return &pb.CreateGroupResponse{GroupId: resGroupId}, nil
}

func (s *CardServer) ListGroups(ctx context.Context, _ *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groups, err := s.GroupService.List(ctx, userId)

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

func (s *CardServer) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	groupId := entity.GroupId(req.GroupId)

	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	group, err := s.GroupService.GetGroup(ctx, userId, groupId)

	if err != nil {
		return nil, err
	}

	groupResp := groupToResponse(group)
	return &pb.GetGroupResponse{Group: groupResp}, nil
}

func (s *CardServer) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*emptypb.Empty, error) {
	groupId := entity.GroupId(req.GroupId)
	visibility := entity.GroupVisibility(req.Visibility)
	group := entity.UpdateGroup{
		Id:          groupId,
		Name:        req.GroupName,
		Description: req.Description,
		Visibility:  visibility,
	}

	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.GroupService.UpdateGroup(ctx, userId, group); err != nil {
		return nil, err
	}

	s.Logger.Info("Group updated", slog.Int64("groupId", req.GroupId))
	return &emptypb.Empty{}, nil
}

func (s *CardServer) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*emptypb.Empty, error) {
	groupId := entity.GroupId(req.GroupId)

	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.GroupService.DeleteGroup(ctx, userId, groupId); err != nil {
		return nil, err
	}

	s.Logger.Info("Group deleted", slog.Int64("groupId", req.GroupId))
	return &emptypb.Empty{}, nil
}

func (s *CardServer) AddCard(ctx context.Context, req *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groupId := entity.GroupId(req.GroupId)
	cardId, err := s.CardsService.Create(ctx, userId, groupId, req.FrontText, req.BackText)

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

func (s *CardServer) ListCards(ctx context.Context, req *pb.ListCardsRequest) (
	*pb.ListCardsResponse,
	error,
) {
	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	groupId := entity.GroupId(req.GroupId)
	cards, err := s.CardsService.ListCards(ctx, userId, groupId)

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

func (s *CardServer) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	cardId := entity.CardId(req.CardId)
	card, err := s.CardsService.GetCard(ctx, userId, cardId)

	if err != nil {
		return nil, err
	}

	cardResp := cardToResponse(card)
	return &pb.GetCardResponse{Card: cardResp}, nil
}

func (s *CardServer) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*emptypb.Empty, error) {
	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	cardId := entity.CardId(req.CardId)
	card := entity.UpdateCard{Id: cardId, FrontText: req.FrontText, BackText: req.BackText}

	if err = s.CardsService.UpdateCard(ctx, userId, card); err != nil {
		return nil, err
	}

	s.Logger.Info("Card updated", slog.Int64("cardId", req.CardId))
	return &emptypb.Empty{}, nil
}

func (s *CardServer) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*emptypb.Empty, error) {
	userId, err := s.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	cardId := entity.CardId(req.CardId)

	if err = s.CardsService.DeleteCard(ctx, userId, cardId); err != nil {
		return nil, err
	}

	s.Logger.Info("Card deleted", slog.Int64("cardId", req.CardId))
	return &emptypb.Empty{}, nil
}

func (s *CardServer) HealthCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
