package grpc

import (
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"time"
)

func groupToResponse(group entity.Group) *pb.CardsGroup {
	return &pb.CardsGroup{
		Id:          int64(group.Id),
		OwnerId:     int64(group.OwnerId),
		Name:        group.Name,
		Description: group.Description,
		CreatedAt:   group.CreateTime.Format(time.RFC3339),
		Visibility:  pb.GroupVisibility(group.Visibility),
	}
}

func cardToResponse(card entity.Card) *pb.Card {
	return &pb.Card{
		Id:        int64(card.Id),
		GroupId:   int64(card.GroupId),
		FrontText: card.FrontText,
		BackText:  card.BackText,
		CreatedAt: card.CreateTime.Format(time.RFC3339),
	}
}
