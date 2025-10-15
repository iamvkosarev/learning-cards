//go:build !mecab

package japanese

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"log/slog"
)

type Reader struct {
	Config config.JapaneseReading
}

func NewReader(config config.JapaneseReading, logger *slog.Logger) *Reader {
	return &Reader{
		Config: config,
	}
}

func (j Reader) GetCardReading(ctx context.Context, text string) ([]model.ReadingPair, error) {
	return []model.ReadingPair{}, nil
}
