package cards

import (
	"context"
	"log"

	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
)

type Handler struct {
	pb.UnimplementedLearningCardsServer
}

func NewCardsHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CreateCardsGroup(
	ctx context.Context,
	req *pb.CreateCardsGroupRequest,
) (*pb.CreateCardsGroupResponse, error) {
	log.Println("CreateCardsGroup: received:", req.GroupName)
	return &pb.CreateCardsGroupResponse{
		GroupId: -1,
	}, nil
}
