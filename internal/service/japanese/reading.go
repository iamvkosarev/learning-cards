//go:build mecab

package japanese

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/shogo82148/go-mecab"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"runtime/debug"
	"strings"
	"time"

	"unicode"
)

const mecabReaderTracerName = "mecab.reader"

var ErrFailedToAnalyzeHiragana = errors.New("failed to analyze hiragana")

type Reader struct {
	Config config.JapaneseReading
	tracer trace.Tracer
	logger *slog.Logger
}

func NewReader(config config.JapaneseReading, logger *slog.Logger) *Reader {
	return &Reader{
		Config: config,
		tracer: otel.Tracer(mecabReaderTracerName),
		logger: logger,
	}
}

type readingPairsResult struct {
	pairs []model.ReadingPair
	err   error
}

func (j *Reader) GetCardReading(ctx context.Context, text string) ([]model.ReadingPair, error) {
	ctx, span := j.tracer.Start(ctx, "GetCardReading")
	defer span.End()

	op := "service.Reader.GetCardReading"
	ch := make(chan *readingPairsResult)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				j.logger.Error(
					"failed to analyze hiragana", slog.String("reason", fmt.Sprintf("%v", r)),
					slog.String("stack", string(debug.Stack())), slog.String("input", text),
				)
				ch <- &readingPairsResult{
					pairs: nil,
					err:   ErrFailedToAnalyzeHiragana,
				}
			}
		}()
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
func (j *Reader) analyzeWithHiragana(text string) ([]model.ReadingPair, error) {
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
	lastReadingIndex := 0

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
		var dicForm string
		if len(features) >= 7 && features[6] != "*" {
			dicForm = features[6]
		} else {
			dicForm = surface
		}
		if !containsKanji(dicForm) {
			lastReadingIndex += len(dicForm)
			continue
		}

		originForm := trimHiraganaFromNotEqualEnd(text[lastReadingIndex:lastReadingIndex+len(dicForm)], reading)
		lastReadingIndex += len(dicForm)

		result = append(
			result, model.ReadingPair{
				Text:    originForm,
				Reading: reading,
			},
		)
	}
	return result, nil
}

func trimHiraganaFromNotEqualEnd(text string, reading string) string {
	textRunes := []rune(text)
	if unicode.In(textRunes[len(textRunes)-1], unicode.Han) {
		return text
	}
	readingRunes := []rune(reading)
	lastIndex := len(textRunes) - 1
	readingRunesIndex := len(readingRunes) - 1
	for i := len(textRunes) - 1; i >= 0; i-- {
		if textRunes[i] != readingRunes[readingRunesIndex] {
			lastIndex = i + 1
		} else {
			readingRunesIndex--
		}
	}
	return string(textRunes[:lastIndex+1])
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
