package discord_client

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/db"
	"golang.org/x/exp/slog"
)

type DiscordClient struct {
	cfg            config.Config
	queries        *db.Queries
	session        *discordgo.Session
	channelIDCache map[string]*discordgo.Channel
	channelsCache  []*discordgo.Channel
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
		channelIDCache: map[string]*discordgo.Channel{},
		channelsCache:  []*discordgo.Channel{},
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

func (d *DiscordClient) MessageReactionRemove(channelID string, messageID string, emoji string) error {
	return d.session.MessageReactionRemove(channelID, messageID, emoji, d.session.State.User.ID)
}

func (d *DiscordClient) GetMessageReactions(channelID string, messageID string, emoji string) ([]*discordgo.User, error) {
	return d.session.MessageReactions(channelID, messageID, emoji, 100, "", "")
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

func (d *DiscordClient) UpdateChannelsCache() error {
	guildID := d.session.State.Guilds[0].ID
	channels, err := d.session.GuildChannels(guildID)
	if err != nil {
		return err
	}
	d.channelsCache = channels
	return nil
}

func (d *DiscordClient) GetCachedChannel(category string, origChannelName string) (*discordgo.Channel, error) {
	category = strings.ToLower(category)
	channel := strings.ToLower(origChannelName)
	channel = strings.ReplaceAll(channel, "\u3000", "-")
	cacheKey := category + "/" + channel
	cachedChannel, ok := d.channelIDCache[cacheKey]
	if ok {
		return cachedChannel, nil
	}

	guildID := d.session.State.Guilds[0].ID
	categoryID := ""
	for _, ch := range d.channelsCache {
		if ch.Type == discordgo.ChannelTypeGuildCategory && ch.Name == category {
			categoryID = ch.ID
			break
		}
	}
	if categoryID == "" {
		categoryChannel, err := d.session.GuildChannelCreate(guildID, category, discordgo.ChannelTypeGuildCategory)
		if err != nil {
			return nil, err
		}
		data := discordgo.GuildChannelCreateData{
			Name:     channel,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: categoryChannel.ID,
		}
		createdChannel, err := d.session.GuildChannelCreateComplex(guildID, data)
		if err != nil {
			return nil, err
		}
		slog.Debug("GuildChannelCreateComplex OK", "name", channel, "cacheKey", cacheKey, "created ch.Name", createdChannel.Name)
		d.channelIDCache[cacheKey] = createdChannel
		return createdChannel, nil
	}
	for _, ch := range d.channelsCache {
		if ch.Type == discordgo.ChannelTypeGuildText && ch.ParentID == categoryID && ch.Name == channel {
			d.channelIDCache[cacheKey] = ch
			return ch, nil
		}
	}
	data := discordgo.GuildChannelCreateData{
		Name:     channel,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	}
	ch, err := d.session.GuildChannelCreateComplex(guildID, data)
	if err != nil {
		return nil, err
	}
	slog.Debug("GuildChannelCreateComplex OK", "origChannelName", origChannelName, "cacheKey", cacheKey, "created ch.Name", ch.Name)
	d.channelIDCache[cacheKey] = ch
	return ch, nil
}

func (d *DiscordClient) SendMessage(category string, channel string, message string) (string, error) {
	if len(d.session.State.Guilds) != 1 {
		return "", fmt.Errorf("discord app must join one server")
	}
	ch, err := d.GetCachedChannel(category, channel)
	if err != nil {
		return "", err
	}
	msg, err := d.session.ChannelMessageSend(ch.ID, message)
	if err != nil {
		return "", err
	} else {
		return msg.ID, nil
	}
}
