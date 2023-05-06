package discord_client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"go.uber.org/zap"
	"golang.org/x/text/width"
)

type DiscordClient struct {
	cfg            config.Config
	queries        *db.Queries
	session        *discordgo.Session
	logger         *zap.Logger
	channelIDCache map[string]*discordgo.Channel
	channelsCache  []*discordgo.Channel
}

func NewDiscordClient(cfg config.Config, queries *db.Queries, logger *zap.Logger) (*DiscordClient, error) {
	session, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		return nil, err
	}
	session.Identify.Intents |= discordgo.IntentsMessageContent
	return &DiscordClient{
		cfg:            cfg,
		queries:        queries,
		session:        session,
		logger:         logger,
		channelIDCache: map[string]*discordgo.Channel{},
		channelsCache:  []*discordgo.Channel{},
	}, nil
}

func (d *DiscordClient) Session() *discordgo.Session {
	return d.session
}

func (d *DiscordClient) GetChannel(channelID string) (*discordgo.Channel, error) {
	return d.session.Channel(channelID)
}

func (d *DiscordClient) GetChannelMessage(channelID string, messageID string) (*discordgo.Message, error) {
	return d.session.ChannelMessage(channelID, messageID)
}

func (d *DiscordClient) MessageReactionAdd(channelID string, messageID string, emoji string) error {
	err := d.session.MessageReactionAdd(channelID, messageID, emoji)
	return err
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
	if len(d.session.State.Guilds) != 1 {
		return errors.New("bot must join exactly one server")
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

func (d *DiscordClient) GetCachedChannel(origCategory string, origChannelName string) (*discordgo.Channel, error) {
	category := strings.ToLower(width.Fold.String(origCategory))
	category = strings.ReplaceAll(category, "\u3000", "-")
	category = strings.ReplaceAll(category, " ", "-")
	channel := strings.ToLower(width.Fold.String(origChannelName))
	channel = strings.ReplaceAll(channel, "\u3000", "-")
	channel = strings.ReplaceAll(channel, " ", "-")
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
		d.channelsCache = append(d.channelsCache, categoryChannel)
		data := discordgo.GuildChannelCreateData{
			Name:     channel,
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: categoryChannel.ID,
		}
		createdChannel, err := d.session.GuildChannelCreateComplex(guildID, data)
		if err != nil {
			return nil, err
		}
		d.channelsCache = append(d.channelsCache, createdChannel)
		d.logger.Debug("GuildChannelCreateComplex OK", zap.String("name", channel), zap.String("cacheKey", cacheKey), zap.String("created ch.Name", createdChannel.Name))
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
	d.channelsCache = append(d.channelsCache, ch)
	d.logger.Debug("GuildChannelCreateComplex OK", zap.String("origChannelName", origChannelName), zap.String("cacheKey", cacheKey), zap.String("created ch.Name", ch.Name))
	d.channelIDCache[cacheKey] = ch
	return ch, nil
}

func (d *DiscordClient) SendMessage(category string, channel string, message string) (*discordgo.Message, error) {
	if len(d.session.State.Guilds) != 1 {
		return nil, fmt.Errorf("discord app must join one server [%d]", len(d.session.State.Guilds))
	}
	ch, err := d.GetCachedChannel(category, channel)
	if err != nil {
		return nil, err
	}
	msg, err := d.session.ChannelMessageSend(ch.ID, message)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (d *DiscordClient) EditMessage(category string, channel string, messageID, message string) (*discordgo.Message, error) {
	if len(d.session.State.Guilds) != 1 {
		return nil, fmt.Errorf("discord app must join one server [%d]", len(d.session.State.Guilds))
	}
	ch, err := d.GetCachedChannel(category, channel)
	if err != nil {
		return nil, err
	}
	msg, err := d.session.ChannelMessageEdit(ch.ID, messageID, message)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (d *DiscordClient) createChannelWithTopic(category string, channel string, topic string) (*discordgo.Channel, error) {
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
			Topic:    topic,
		}
		createdChannel, err := d.session.GuildChannelCreateComplex(guildID, data)
		if err != nil {
			return nil, err
		}
		d.logger.Debug("GuildChannelCreateComplex OK", zap.String("name", channel), zap.String("created ch.Name", createdChannel.Name))
		return createdChannel, nil
	}
	for _, ch := range d.channelsCache {
		if ch.Type == discordgo.ChannelTypeGuildText && ch.ParentID == categoryID && ch.Name == channel {
			edit := discordgo.ChannelEdit{
				Topic: topic,
			}
			d.session.ChannelEdit(ch.ID, &edit)
			return ch, nil
		}
	}
	data := discordgo.GuildChannelCreateData{
		Name:     channel,
		Type:     discordgo.ChannelTypeGuildText,
		Topic:    topic,
		ParentID: categoryID,
	}
	ch, err := d.session.GuildChannelCreateComplex(guildID, data)
	if err != nil {
		return nil, err
	}
	d.logger.Debug("GuildChannelCreateComplex OK", zap.String("created ch.Name", ch.Name))
	return ch, nil
}

func (d *DiscordClient) CreateNotifyAndScheduleChannel() (*discordgo.Channel, error) {
	ch, err := d.createChannelWithTopic(discord.NotifyAndScheduleCategory, discord.AutoSearchChannelName, discord.AutoSearchChannelTopic)
	if err != nil {
		return nil, err
	}
	err = d.UpdateChannelsCache()
	if err != nil {
		return nil, err
	}
	msgs, err := d.session.ChannelMessages(ch.ID, 1, "", "0", "")
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		_, err := d.SendMessage(discord.NotifyAndScheduleCategory, discord.AutoSearchChannelName, discord.AutoSearchChannelWelcomeMessage)
		if err != nil {
			return nil, err
		}
	} else {
		if msgs[0].Author.ID == d.session.State.User.ID {
			_, err := d.session.ChannelMessageEdit(ch.ID, msgs[0].ID, discord.AutoSearchChannelWelcomeMessage)
			if err != nil {
				return nil, err
			}
		}
	}
	return ch, err
}

func (d *DiscordClient) SendMessageToThread(threadID string, content string) error {
	_, err := d.session.ChannelMessageSend(threadID, content)
	return err
}

func (d *DiscordClient) CreateAutoSearchThread(title string, content string) (string, error) {
	ch, err := d.GetCachedChannel(discord.NotifyAndScheduleCategory, discord.AutoSearchChannelName)
	if err != nil {
		return "", err
	}
	th, err := d.session.ThreadStart(ch.ID, title, discordgo.ChannelTypeGuildPublicThread, 7*24*60)
	if err != nil {
		return "", err
	}
	_, err = d.session.ChannelMessageSend(th.ID, content)
	return th.ID, nil
}

func (d *DiscordClient) DeleteThread(threadID string) error {
	_, err := d.session.ChannelDelete(threadID)
	return err
}
