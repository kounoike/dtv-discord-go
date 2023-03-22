package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/db"
)

var (
	recordingReactionEmoji = "📼"
	okReactionEmoji        = "🆗"
)

type DiscordClient struct {
	ctx            context.Context
	cfg            config.Config
	queries        *db.Queries
	session        *discordgo.Session
	channelIDCache map[string]string
}

func NewDiscordClient(ctx context.Context, cfg config.Config, queries *db.Queries) (*DiscordClient, error) {
	session, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		return nil, err
	}
	return &DiscordClient{
		ctx:            ctx,
		cfg:            cfg,
		queries:        queries,
		session:        session,
		channelIDCache: map[string]string{},
	}, nil
}

func (d *DiscordClient) reactionAdd(discord *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	fmt.Println("add", reaction.Emoji.Name, reaction.UserID, reaction.ChannelID)

	if reaction.Emoji.Name != recordingReactionEmoji {
		return
	}

	msg, err := discord.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range msg.Reactions {
		if r.Emoji.Name == recordingReactionEmoji {
			if r.Count == 1 {
				// 録画しよう！
				programMessage, err := d.queries.GetProgramMessageByMessageID(d.ctx, reaction.MessageID)
				if err != nil {
					fmt.Println(err)
					return
				}
				url := fmt.Sprintf("http://%s:%d/api/recording/schedules", d.cfg.Mirakc.Host, d.cfg.Mirakc.Port)
				postOption := fmt.Sprintf(`{"programId": %d, "options": {"contentPath": "%d.m2ts"}, "tags": ["manual"]}`, programMessage.ProgramID, programMessage.ProgramID)
				client := resty.New()
				resp, err := client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(postOption).
					Post(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("録画予約: StatusCode:", resp.StatusCode())
				if resp.StatusCode() == 201 {
					// 録画予約成功？
					err := d.session.MessageReactionAdd(reaction.ChannelID, reaction.MessageID, okReactionEmoji)
					fmt.Println(`録画予約成功絵文字付加`, err)
				}
			}
			break
		}
	}
}

func (d *DiscordClient) reactionRemove(discord *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	fmt.Println("remove", reaction.Emoji.Name, reaction.UserID, reaction.ChannelID)
	msg, err := discord.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range msg.Reactions {
		fmt.Println(r.Emoji.Name, r.Count)
		// if r.Emoji.Name == recordingReactionEmoji {
		// 	break
		// }
	}
}

func (d *DiscordClient) Open() error {
	err := d.session.Open()
	if err != nil {
		return err
	}
	d.session.AddHandler(d.reactionAdd)
	d.session.AddHandler(d.reactionRemove)
	return nil
}

func (d *DiscordClient) getChannelID(category string, channel string) (string, error) {
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
	chID, err := d.getChannelID(category, channel)
	if err != nil {
		return "", err
	}
	msg, err := d.session.ChannelMessageSend(chID, message)
	if err != nil {
		return "", err
	} else {
		// fmt.Println(category, channel, chID, msg)
		return msg.ID, nil
	}
}
