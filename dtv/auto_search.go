package dtv

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"go.uber.org/zap"
	"golang.org/x/text/width"
	"gopkg.in/ini.v1"
)

type AutoSearch struct {
	Title          string            `yaml:"タイトル"`
	Channel        string            `yaml:"チャンネル"`
	Genre          string            `yaml:"ジャンル"`
	NotifyUsers    []*discordgo.User `yaml:"-"`
	RecordingUsers []*discordgo.User `yaml:"-"`
	ThreadID       string            `yaml:"-"`
}

type AutoSearchProgram struct {
	Title string
	Genre string
}

func NewAutoSearchProgram(p db.Program, kanaMatch bool) *AutoSearchProgram {
	return &AutoSearchProgram{
		Title: normalizeString(p.Name, kanaMatch),
		Genre: normalizeString(p.Genre, kanaMatch),
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
	if a.Title != "" && !strings.Contains(program.Title, a.Title) {
		return false
	}
	if a.Genre != "" && !strings.Contains(program.Genre, a.Genre) {
		return false
	}
	return true
}

func (a *AutoSearch) IsMatchService(serviceName string, kanaMatch bool) bool {
	if a.Channel == "" || strings.Contains(normalizeString(serviceName, kanaMatch), normalizeString(a.Channel, kanaMatch)) {
		return true
	} else {
		return false
	}
}

func (dtv *DTVUsecase) getAutoSeachFromMessage(msg *discordgo.Message) (*AutoSearch, error) {
	content := []byte(msg.Content)
	iniContent, err := ini.Load(content)
	if err != nil {
		return nil, err
	}
	autoSearch := AutoSearch{}
	autoSearch.Channel = normalizeString(iniContent.Section("").Key("チャンネル").String(), dtv.kanaMatch)
	autoSearch.Genre = normalizeString(iniContent.Section("").Key("ジャンル").String(), dtv.kanaMatch)
	autoSearch.Title = normalizeString(iniContent.Section("").Key("タイトル").String(), dtv.kanaMatch)

	notifyUsers, err := dtv.discord.GetMessageReactions(msg.ChannelID, msg.ID, discord.NotifyReactionEmoji)
	if err != nil {
		dtv.logger.Warn("can't get message reactions", zap.Error(err), zap.String("msg.ChannelID", msg.ChannelID), zap.String("msg.ID", msg.ID), zap.String("emoji", discord.NotifyReactionEmoji))
		notifyUsers = []*discordgo.User{}
	}
	autoSearch.NotifyUsers = notifyUsers

	recordingUsers, err := dtv.discord.GetMessageReactions(msg.ChannelID, msg.ID, discord.RecordingReactionEmoji)
	if err != nil {
		dtv.logger.Warn("can't get message reactions", zap.Error(err), zap.String("msg.ChannelID", msg.ChannelID), zap.String("msg.ID", msg.ID), zap.String("emoji", discord.RecordingReactionEmoji))
		recordingUsers = []*discordgo.User{}
	}
	autoSearch.RecordingUsers = recordingUsers
	autoSearch.ThreadID = msg.ChannelID

	return &autoSearch, nil
}

func (dtv *DTVUsecase) ListAutoSearchForServiceName(serviceName string) ([]*AutoSearch, error) {
	msgs, err := dtv.discord.ListAutoSearchChannelThredOkReactionedFirstMessageContents(dtv.autoSearchChannel.ID)
	if err != nil {
		return nil, err
	}
	serviceNameNormalized := normalizeString(serviceName, dtv.kanaMatch)

	autoSearchList := make([]*AutoSearch, 0)

	for _, msg := range msgs {
		autoSearch, err := dtv.getAutoSeachFromMessage(msg)
		if err != nil {
			dtv.logger.Warn("thread message yaml unmarshal error", zap.Error(err))
			continue
		}
		if autoSearch.Channel == "" || strings.Contains(serviceNameNormalized, normalizeString(autoSearch.Channel, dtv.kanaMatch)) {
			autoSearch.Title = normalizeString(autoSearch.Title, dtv.kanaMatch)
			autoSearchList = append(autoSearchList, autoSearch)
		}
	}
	return autoSearchList, nil
}
