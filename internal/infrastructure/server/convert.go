package server

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

func cardToReviewResponse(card entity.Card) *pb.ReviewCard {
	return &pb.ReviewCard{
		Id:        int64(card.Id),
		FrontText: card.FrontText,
		BackText:  card.BackText,
	}
}
func markToResponse(mark entity.Mark) pb.Mark {
	switch mark {
	case entity.MARK_A:
		return pb.Mark_MARK_A
	case entity.MARK_B:
		return pb.Mark_MARK_B
	case entity.MARK_C:
		return pb.Mark_MARK_C
	case entity.MARK_D:
		return pb.Mark_MARK_D
	case entity.MARK_E:
		return pb.Mark_MARK_E
	default:
		return pb.Mark_MARK_NULL
	}
}

func answerToEntity(answer pb.CardAnswer) entity.Answer {
	switch answer {
	case pb.CardAnswer_EASY:
		return entity.ANSWER_EASY
	case pb.CardAnswer_FAIL:
		return entity.ANSWER_FAIL
	case pb.CardAnswer_HARD:
		return entity.ANSWER_HARD
	case pb.CardAnswer_GOOD:
		return entity.ANSWER_GOOD
	}
	return entity.ANSWER_EASY
}
