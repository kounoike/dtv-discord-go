package dtv

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"go.uber.org/zap"
	"golang.org/x/text/width"
	"gopkg.in/yaml.v3"
)

type AutoSearch struct {
	Title          string            `yaml:"タイトル"`
	Channel        string            `yaml:"チャンネル"`
	NotifyUsers    []*discordgo.User `yaml:"-"`
	RecordingUsers []*discordgo.User `yaml:"-"`
	ThreadID       string            `yaml:"-"`
}

type AutoSearchProgram struct {
	Title string
}

func NewAutoSearchProgram(p db.Program) *AutoSearchProgram {
	return &AutoSearchProgram{
		Title: normalizeString(p.Name),
	}
}

func normalizeString(str string) string {
	return strings.ToLower(width.Fold.String(str))
}

func (a *AutoSearch) IsMatchProgram(program *AutoSearchProgram) bool {
	if strings.Contains(program.Title, a.Title) {
		return true
	} else {
		return false
	}
}

func (dtv *DTVUsecase) ListAutoSearchForServiceName(serviceName string) ([]*AutoSearch, error) {
	msgs, err := dtv.discord.ListAutoSearchChannelThredFirstMessageContents(dtv.autoSearchChannel.ID)
	if err != nil {
		return nil, err
	}
	serviceNameNormalized := normalizeString(serviceName)

	autoSearchList := make([]*AutoSearch, 0)

	for _, msg := range msgs {
		content := []byte(msg.Content)
		var autoSearch AutoSearch
		err := yaml.Unmarshal(content, &autoSearch)
		if err != nil {
			dtv.logger.Warn("thread message yaml unmarshal error", zap.Error(err))
			continue
		}
		if autoSearch.Channel == "" || strings.Contains(serviceNameNormalized, normalizeString(autoSearch.Channel)) {
			autoSearch.Title = normalizeString(autoSearch.Title)

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

			autoSearchList = append(autoSearchList, &autoSearch)
		}
	}
	return autoSearchList, nil
}
