package service

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/shogo82148/go-mecab"
	"strings"
	"time"
	"unicode"
)

type JapaneseReader struct {
	Config config.JapaneseReading
}

func NewJapaneseReader(config config.JapaneseReading) *JapaneseReader {
	return &JapaneseReader{
		Config: config,
	}
}

type readingPairsResult struct {
	pairs []model.ReadingPair
	err   error
}

func (j *JapaneseReader) GetCardReading(ctx context.Context, text string) ([]model.ReadingPair, error) {
	op := "service.JapaneseReader.GetCardReading"
	ch := make(chan *readingPairsResult)
	go func() {
		pairs, err := j.analyzeWithHiragana(text)
		ch <- &readingPairsResult{
			pairs: pairs,
			err:   err,
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-ch:
		return result.pairs, result.err
	case <-time.After(j.Config.SearchTimeout):
		return nil, fmt.Errorf("%s: %w", op, model.ErrTimeOut)
	}
	return nil, nil
}

func (j *JapaneseReader) analyzeWithHiragana(text string) ([]model.ReadingPair, error) {
	mecabModel, err := mecab.NewModel(
		map[string]string{
			"dicdir": j.Config.MecabDicDir,
		},
	)
	if err != nil {
		return nil, err
	}
	defer mecabModel.Destroy()

	tagger, err := mecabModel.NewMeCab()
	if err != nil {
		return nil, err
	}
	defer tagger.Destroy()

	node, err := tagger.ParseToNode(text)
	if err != nil {
		return nil, err
	}

	var result []model.ReadingPair
	for ; !node.IsZero(); node = node.Next() {
		if node.Stat() == mecab.BOSNode || node.Stat() == mecab.EOSNode {
			continue
		}

		surface := node.Surface()
		features := strings.Split(node.Feature(), ",")
		var reading string
		if len(features) >= 8 && features[7] != "*" {
			reading = KatakanaToHiragana(features[7])
		} else {
			reading = surface
		}
		var base string
		if len(features) >= 7 && features[6] != "*" {
			base = features[6]
		} else {
			base = surface
		}
		if !containsKanji(base) {
			continue
		}

		result = append(
			result, model.ReadingPair{
				Text:    base,
				Reading: reading,
			},
		)
	}
	return result, nil
}

func KatakanaToHiragana(input string) string {
	var result strings.Builder
	for _, r := range input {
		if r >= 'ァ' && r <= 'ン' {
			r -= 0x60
		}
		result.WriteRune(r)
	}
	return result.String()
}

func containsKanji(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han) {
			return true
		}
	}
	return false
}
