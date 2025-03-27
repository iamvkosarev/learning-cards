package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

const internalErrMessage = "internal server error"

type groupUseCase interface {
	Create(ctx context.Context, name, description string, visibility entity.GroupVisibility) (entity.GroupId, error)
	List(ctx context.Context) ([]entity.Group, error)
	Get(ctx context.Context, id entity.GroupId) (entity.Group, error)
	Update(ctx context.Context, updateGroup entity.UpdateGroup) error
}

type cardsUseCase interface {
	Create(ctx context.Context, groupId entity.GroupId, frontText, backText string) (entity.CardId, error)
	List(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error)
	Get(ctx context.Context, id entity.CardId) (entity.Card, error)
	Update(ctx context.Context, card entity.UpdateCard) error
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

func (s *Server) CreateCardsGroup(ctx context.Context, req *pb.CreateCardsGroupRequest) (
	*pb.CreateCardsGroupResponse,
	error,
) {
	const op = "grpc.CreateCardsGroup"
	log := s.Logger.With(slog.String("op", op))

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	visibility := entity.GroupVisibility(req.Visibility)
	cardId, err := s.groupUseCase.Create(ctx, req.GroupName, req.Description, visibility)

	if verificationErr := getVerificationErr(log, err); verificationErr != nil {
		return nil, verificationErr
	}

	if err != nil {
		log.Info("failed to create group", sl.Err(err))
		return nil, status.Error(codes.Internal, internalErrMessage)
	}

	resId := int64(cardId)
	log.Info("Cards group created", slog.String("name", req.GroupName), slog.Int64("id", resId))
	return &pb.CreateCardsGroupResponse{GroupId: resId}, nil
}

func (s *Server) ListCardsGroups(ctx context.Context, _ *emptypb.Empty) (*pb.ListCardsGroupsResponse, error) {
	const op = "grpc.ListCardsGroups"
	log := s.Logger.With(slog.String("op", op))

	groups, err := s.groupUseCase.List(ctx)

	if verificationErr := getVerificationErr(log, err); verificationErr != nil {
		return nil, verificationErr
	}
	if err != nil {
		log.Info("failed to create card", sl.Err(err))
		return nil, status.Error(codes.Internal, internalErrMessage)
	}

	var respGroups []*pb.CardsGroup
	for _, group := range groups {
		respGroups = append(
			respGroups, groupToResponse(group),
		)
	}

	return &pb.ListCardsGroupsResponse{Groups: respGroups}, nil
}

func (s *Server) GetCardsGroupCards(ctx context.Context, req *pb.GetCardsGroupCardsRequest) (
	*pb.GetCardsGroupCardsResponse,
	error,
) {
	const op = "grpc.GetCardsGroupCards"
	log := s.Logger.With(slog.String("op", op))

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	groupId := entity.GroupId(req.GroupId)
	cards, err := s.cardsUseCase.List(ctx, groupId)

	if verificationErr := getVerificationErr(log, err); verificationErr != nil {
		return nil, verificationErr
	}
	if err != nil {
		log.Info("failed to create card", sl.Err(err))
		return nil, status.Error(codes.Internal, internalErrMessage)
	}

	var respCards []*pb.Card
	for _, card := range cards {
		respCards = append(
			respCards, cardToResponse(card),
		)
	}

	return &pb.GetCardsGroupCardsResponse{Cards: respCards}, nil
}

func (s *Server) GetCardsGroup(ctx context.Context, req *pb.GetCardsGroupRequest) (*pb.GetCardsGroupResponse, error) {
	const op = "grpc.GetCardsGroup"
	log := s.Logger.With(slog.String("op", op))

	groupId := entity.GroupId(req.GroupId)

	group, err := s.groupUseCase.Get(ctx, groupId)

	if verificationErr := getVerificationErr(log, err); verificationErr != nil {
		return nil, verificationErr
	}

	if err != nil {
		log.Info("failed to create card", sl.Err(err))
		return nil, status.Error(codes.Internal, internalErrMessage)
	}

	groupResp := groupToResponse(group)
	return &pb.GetCardsGroupResponse{Group: groupResp}, nil
}

func (s *Server) UpdateCardsGroup(ctx context.Context, req *pb.UpdateCardsGroupRequest) (
	*emptypb.Empty,
	error,
) {
	const op = "grpc.UpdateCardsGroup"
	log := s.Logger.With(
		slog.String("op", op),
		slog.Int64("groupId", req.GroupId),
	)

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	groupId := entity.GroupId(req.GroupId)
	visibility := entity.GroupVisibility(req.Visibility)

	err = s.groupUseCase.Update(
		ctx, entity.UpdateGroup{Id: groupId, Name: req.GroupName, Description: req.Description, Visibility: visibility},
	)

	if verificationErr := getVerificationErr(log, err); verificationErr != nil {
		return nil, verificationErr
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrGroupNotFound):
			return nil, status.Error(codes.NotFound, "group not found")
		case errors.Is(err, entity.ErrCardNotFound):
			return nil, status.Error(codes.NotFound, "card not found")
		default:
			log.Info("failed to update group", sl.Err(err))
			return nil, status.Error(codes.Internal, internalErrMessage)
		}
	}

	log.Info("Group updated")
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCardsGroup(ctx context.Context, req *pb.DeleteCardsGroupRequest) (*emptypb.Empty, error) {
	const op = "grpc.DeleteCardsGroup"
	//log := s.Logger.With(slog.String("op", op))
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) AddCard(ctx context.Context, req *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	const op = "grpc.AddCard"
	log := s.Logger.With(
		slog.String("op", op),
		slog.Int64("groupId", req.GroupId),
	)

	err := req.Validate()
	if err != nil {
		log.Warn("invalid request", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	groupId := entity.GroupId(req.GroupId)
	cardId, err := s.cardsUseCase.Create(ctx, groupId, req.FrontText, req.BackText)

	if verificationErr := getVerificationErr(log, err); verificationErr != nil {
		return nil, verificationErr
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrGroupNotFound):
			return nil, status.Error(
				codes.NotFound,
				fmt.Sprintf("group (id:%v) not found", req.GroupId),
			)
		default:
			s.Logger.Info("failed to create card", sl.Err(err))
			return nil, status.Error(
				codes.Internal,
				internalErrMessage,
			)
		}
	}

	resId := int64(cardId)
	log.Info("Card created", slog.Int64("cardId", resId))
	return &pb.AddCardResponse{CardId: resId}, nil
}

func (s *Server) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	const op = "grpc.GetCard"
	log := s.Logger.With(slog.String("op", op))

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cardId := entity.CardId(req.CardId)
	card, err := s.cardsUseCase.Get(ctx, cardId)

	if verificationErr := getVerificationErr(log, err); verificationErr != nil {
		return nil, verificationErr
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrCardNotFound):
			log.Warn("card not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "card not found")
		default:
			log.Info("failed to get card", sl.Err(err))
			return nil, status.Error(codes.Internal, internalErrMessage)
		}
	}

	cardResp := cardToResponse(card)

	return &pb.GetCardResponse{Card: cardResp}, nil
}

func (s *Server) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*emptypb.Empty, error) {
	const op = "grpc.UpdateCard"
	log := s.Logger.With(
		slog.String("op", op),
		slog.Int64("cardId", req.CardId),
	)

	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Server) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*emptypb.Empty, error) {
	const op = "grpc.DeleteCard"
	//log := s.Logger.With(slog.String("op", op))
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func getVerificationErr(log *slog.Logger, err error) error {
	var verificationErr entity.VerificationError
	switch {
	case errors.Is(err, entity.ErrMetadataIsEmpty):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.As(err, &verificationErr):
		if verificationErr.StatusCode != codes.PermissionDenied &&
			verificationErr.StatusCode != codes.InvalidArgument {
			log.Error("failed to verify user", sl.Err(verificationErr))
			return status.Error(codes.Internal, "internal verification error")
		} else {
			return status.Error(verificationErr.StatusCode, verificationErr.Message)
		}
	}
	return nil
}
