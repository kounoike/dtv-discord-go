package discord_client

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/db"
)

type DiscordClient struct {
	cfg            config.Config
	queries        *db.Queries
	session        *discordgo.Session
	channelIDCache map[string]string
}

func NewDiscordClient(cfg config.Config, queries *db.Queries) (*DiscordClient, error) {
	session, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		return nil, err
	}
	return &DiscordClient{
		cfg:            cfg,
		queries:        queries,
		session:        session,
		channelIDCache: map[string]string{},
	}, nil
}

func (d *DiscordClient) Session() *discordgo.Session {
	return d.session
}

func (d *DiscordClient) GetChannelMessage(channelID string, messageID string) (*discordgo.Message, error) {
	return d.session.ChannelMessage(channelID, messageID)
}

func (d *DiscordClient) MessageReactionAdd(channelID string, messageID string, emoji string) error {
	return d.session.MessageReactionAdd(channelID, messageID, emoji)
}

func (d *DiscordClient) AddHandler(handler interface{}) {
	d.session.AddHandler(handler)
}

func (d *DiscordClient) Open() error {
	err := d.session.Open()
	if err != nil {
		return err
	}
	return nil
}

func (d *DiscordClient) GetChannelID(category string, channel string) (string, error) {
	category = strings.ToLower(category)
	channel = strings.ToLower(channel)
	cacheKey := category + "/" + channel
	id, ok := d.channelIDCache[cacheKey]
	if ok {
		return id, nil
	}

	guildID := d.session.State.Guilds[0].ID
	channels, err := d.session.GuildChannels(guildID)
	if err != nil {
		return "", err
	}
	categoryID := ""
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory && ch.Name == category {
			categoryID = ch.ID
			break
		}
	}
	if categoryID == "" {
		categoryChannel, err := d.session.GuildChannelCreate(guildID, category, discordgo.ChannelTypeGuildCategory)
		if err != nil {
			return "", err
		}
		data := discordgo.GuildChannelCreateData{
			Name:     channel,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: categoryChannel.ID,
		}
		createdChannel, err := d.session.GuildChannelCreateComplex(guildID, data)
		if err != nil {
			return "", err
		}
		d.channelIDCache[cacheKey] = createdChannel.ID
		return createdChannel.ID, nil
	}
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildText && ch.ParentID == categoryID && ch.Name == channel {
			d.channelIDCache[cacheKey] = ch.ID
			return ch.ID, nil
		}
	}
	data := discordgo.GuildChannelCreateData{
		Name:     channel,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	}
	ch, err := d.session.GuildChannelCreateComplex(guildID, data)
	if err != nil {
		return "", err
	}
	d.channelIDCache[cacheKey] = ch.ID
	return ch.ID, nil
}

func (d *DiscordClient) SendMessage(category string, channel string, message string) (string, error) {
	if len(d.session.State.Guilds) != 1 {
		return "", fmt.Errorf("discord app must join one server")
	}
	chID, err := d.GetChannelID(category, channel)
	if err != nil {
		return "", err
	}
	msg, err := d.session.ChannelMessageSend(chID, message)
	if err != nil {
		return "", err
	} else {
		return msg.ID, nil
	}
}
