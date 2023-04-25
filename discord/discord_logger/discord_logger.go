package discord_logger

import (
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/discord/discord_client"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type WithDiscordEncoder struct {
	zapcore.Encoder
	consoleEncoder zapcore.Encoder
	discordEncoder zapcore.Encoder
	discord        *discord_client.DiscordClient
}

func NewWithDiscordEncoder(level zap.AtomicLevel, discord *discord_client.DiscordClient) *WithDiscordEncoder {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
	encoderCfg.EncodeName = zapcore.FullNameEncoder

	discordEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	return &WithDiscordEncoder{
		consoleEncoder: consoleEncoder,
		discordEncoder: discordEncoder,
		discord:        discord,
	}
}

func (e *WithDiscordEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if entry.Level == zap.ErrorLevel || entry.Level == zap.WarnLevel {
		buf, _ := e.discordEncoder.EncodeEntry(entry, fields)
		ch, _ := e.discord.GetCachedChannel(discord.InformationCategory, discord.ErrorLogChannel)
		_, _ = e.discord.Session().ChannelMessageSend(ch.ID, buf.String())
	}
	if entry.Level == zap.InfoLevel {
		buf, _ := e.discordEncoder.EncodeEntry(entry, fields)
		ch, _ := e.discord.GetCachedChannel(discord.InformationCategory, discord.LogChannel)
		_, _ = e.discord.Session().ChannelMessageSend(ch.ID, buf.String())
	}
	return e.consoleEncoder.EncodeEntry(entry, fields)
}
