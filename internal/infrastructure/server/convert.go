package server

import (
	"github.com/iamvkosarev/learning-cards/internal/model"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"time"
)

func groupToResponse(group *model.Group) *pb.CardsGroup {
	return &pb.CardsGroup{
		Id:          int64(group.Id),
		OwnerId:     int64(group.OwnerId),
		Name:        group.Name,
		Description: group.Description,
		CreatedAt:   group.CreateTime.Format(time.RFC3339),
		Visibility:  pb.GroupVisibility(group.Visibility),
		CardSideTypes: []pb.CardSideType{
			pb.CardSideType(group.CardSideTypes[model.SIDE_FIRST]),
			pb.CardSideType(group.CardSideTypes[model.SIDE_SECOND]),
		},
	}
}

func cardToResponse(card *model.Card) *pb.Card {
	return &pb.Card{
		Id:      int64(card.Id),
		GroupId: int64(card.GroupId),
		Sides: []*pb.CardSide{
			cardSideToResponse(card.Sides[model.SIDE_FIRST]),
			cardSideToResponse(card.Sides[model.SIDE_SECOND]),
		},
		CreatedAt: card.CreateTime.Format(time.RFC3339),
	}
}

func cardSideToResponse(card model.CardSide) *pb.CardSide {
	readingPairs := make([]*pb.ReadingPair, len(card.ReadingPairs))
	for i, readingPair := range card.ReadingPairs {
		readingPairs[i] = &pb.ReadingPair{
			Reading: readingPair.Reading,
			Text:    readingPair.Text,
		}
	}
	return &pb.CardSide{
		Text:         card.Text,
		ReadingPairs: readingPairs,
	}
}

func markToResponse(mark model.Mark) pb.Mark {
	switch mark {
	case model.MARK_A:
		return pb.Mark_MARK_A
	case model.MARK_B:
		return pb.Mark_MARK_B
	case model.MARK_C:
		return pb.Mark_MARK_C
	case model.MARK_D:
		return pb.Mark_MARK_D
	case model.MARK_E:
		return pb.Mark_MARK_E
	default:
		return pb.Mark_MARK_NULL
	}
}

func cardSideTypesToModel(cardSideTypes []pb.CardSideType) []model.CardSideType {
	types := make([]model.CardSideType, len(cardSideTypes))
	for i, sideType := range cardSideTypes {
		switch sideType {
		case pb.CardSideType_CARD_SIDE_JAPANESE:
			types[i] = model.CARD_SIDE_TYPE_JAPANESE
		default:
			types[i] = model.CARD_SIDE_TYPE_NULL
		}
	}
	return types
}

func answerToModels(answer pb.CardAnswer) model.Answer {
	switch answer {
	case pb.CardAnswer_EASY:
		return model.ANSWER_EASY
	case pb.CardAnswer_FAIL:
		return model.ANSWER_FAIL
	case pb.CardAnswer_HARD:
		return model.ANSWER_HARD
	case pb.CardAnswer_GOOD:
		return model.ANSWER_GOOD
	}
	return model.ANSWER_EASY
}
