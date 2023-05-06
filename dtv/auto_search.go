package dtv

import (
	"context"
	"strings"
	"time"

	"github.com/ikawaha/kagome/tokenizer"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"golang.org/x/text/width"
)

type AutoSearch struct {
	Title      string
	Channel    string
	Genre      string
	KanaMatch  bool
	FuzzyMatch bool
	RegexMatch bool
	ThreadID   string
}

type AutoSearchProgram struct {
	Title     string
	TitleKana string
	Genre     string
	GenreKana string
	EndAt     int64
}

func NewAutoSearchProgram(p db.Program) *AutoSearchProgram {
	return &AutoSearchProgram{
		Title:     normalizeString(p.Name, false),
		TitleKana: normalizeString(p.Name, true),
		Genre:     normalizeString(p.Genre, false),
		GenreKana: normalizeString(p.Genre, true),
		EndAt:     p.StartAt + int64(p.Duration),
	}
}

func normalizeString(str string, kanaMatch bool) string {
	normalized := strings.ToLower(width.Fold.String(str))
	if kanaMatch {
		retKana := ""
		t := tokenizer.New()
		tokens := t.Tokenize(normalized)
		normalizedRune := []rune(normalized)
		for _, token := range tokens {
			if len(token.Features()) > 7 {
				retKana += token.Features()[7]
			} else {
				retKana += string(normalizedRune[token.Start:token.End])
			}
		}
		return retKana
	} else {
		return normalized
	}
}

func (a *AutoSearch) IsMatchProgram(program *AutoSearchProgram) bool {
	if program.EndAt < time.Now().Unix() {
		return false
	}
	pTitle := program.Title
	pGenre := program.Genre
	if a.KanaMatch {
		pTitle = program.TitleKana
		pGenre = program.GenreKana
	}
	if a.FuzzyMatch {
		if a.Title != "" && !fuzzy.Match(a.Title, pTitle) {
			return false
		}
		if a.Genre != "" && !fuzzy.Match(a.Genre, pGenre) {
			return false
		}
		return true
	} else {
		if a.Title != "" && !strings.Contains(pTitle, a.Title) {
			return false
		}
		if a.Genre != "" && !strings.Contains(pGenre, a.Genre) {
			return false
		}
		return true
	}
}

func (a *AutoSearch) IsMatchService(serviceName string) bool {
	// NOTE: サービスでfuzzyMatchすると「BS11」が「BSフジ・181」にマッチしてしまうので、fuzzyMatchは使わない
	if a.Channel == "" || strings.Contains(normalizeString(serviceName, a.KanaMatch), normalizeString(a.Channel, a.KanaMatch)) {
		return true
	} else {
		return false
	}
}

func (dtv *DTVUsecase) getAutoSearchFromDB(as *db.AutoSearch) *AutoSearch {
	return &AutoSearch{
		Title:      normalizeString(as.Title, as.KanaSearch),
		Channel:    as.Channel,
		Genre:      normalizeString(as.Genre, as.KanaSearch),
		KanaMatch:  as.KanaSearch,
		FuzzyMatch: as.FuzzySearch,
		RegexMatch: as.RegexSearch,
		ThreadID:   as.ThreadID,
	}
}

func (dtv *DTVUsecase) ListAutoSearchForServiceName(ctx context.Context, serviceName string) ([]*AutoSearch, error) {
	asList, err := dtv.queries.ListAutoSearch(ctx)
	if err != nil {
		return nil, err
	}
	autoSearchList := make([]*AutoSearch, 0, len(asList))

	for _, as := range asList {
		autoSearch := dtv.getAutoSearchFromDB(&as)
		if autoSearch.IsMatchService(serviceName) {
			autoSearchList = append(autoSearchList, autoSearch)
		}
	}
	return autoSearchList, nil
}
