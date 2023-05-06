package dtv

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/template"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/width"
)

func (dtv *DTVUsecase) onProgramsUpdated(ctx context.Context, serviceId uint) error {
	service, err := dtv.mirakc.GetService(serviceId)
	if err != nil {
		return err
	}
	dtv.logger.Debug("Start onProgramsUpdated", zap.Uint("serviceID", serviceId), zap.String("serviceName", service.Name))
	_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.LogChannel, fmt.Sprintf("programs updated: %s", service.Name))
	if err != nil {
		return err
	}
	programs, err := dtv.mirakc.ListPrograms(serviceId)
	if err != nil {
		return err
	}

	autoSearchList, err := dtv.ListAutoSearchForServiceName(ctx, service.Name)
	if err != nil {
		return err
	}

	for _, p := range programs {
		if p.Name == "" {
			continue
		}
		program, err := dtv.queries.GetProgram(ctx, p.ID)
		if errors.Cause(err) == sql.ErrNoRows {
			content, err := template.GetProgramMessage(p, *service)
			if err != nil {
				return err
			}
			content = width.Fold.String(content)
			msg, err := dtv.discord.SendMessage(discord.ProgramInformationCategory, service.Name, content)
			if err != nil {
				return err
			}
			err = dtv.discord.MessageReactionAdd(msg.ChannelID, msg.ID, discord.RecordingReactionEmoji)
			if err != nil {
				return err
			}
			dtv.logger.Debug("will insert program", zap.String("p.Genre", p.Genre))
			err = dtv.queries.InsertProgram(ctx, p)
			if err != nil {
				return err
			}
			err = dtv.queries.InsertProgramMessage(ctx, db.InsertProgramMessageParams{MessageID: msg.ID, ProgramID: p.ID, ChannelID: msg.ChannelID})
			if err != nil {
				return err
			}
			params := db.InsertProgramServiceParams{
				ProgramID: p.ID,
				ServiceID: service.ID,
			}
			err = dtv.queries.InsertProgramService(ctx, params)
			if err != nil {
				return err
			}
			asp := NewAutoSearchProgram(p)
			for _, as := range autoSearchList {
				if as.IsMatchProgram(asp) {
					dtv.logger.Debug("program matched", zap.String("program.Name", p.Name), zap.String("as.Title", as.Title))
					err := dtv.sendAutoSearchMatchMessage(ctx, msg, p, service, as)
					if err != nil {
						dtv.logger.Warn("sendAutoSearchMatchMessage error", zap.Error(err))
						continue
					}
				}
			}
		} else if err != nil {
			return err
		} else {
			if p.Json != program.Json {
				// JSONに違いがあった場合
				dtv.logger.Debug("EPG info updated")

				content, err := template.GetProgramMessage(p, *service)
				if err != nil {
					return err
				}
				content = width.Fold.String(content)
				programMessage, err := dtv.queries.GetProgramMessageByProgramID(ctx, p.ID)
				if err != nil {
					return err
				}
				msg, err := dtv.discord.EditMessage(discord.ProgramInformationCategory, service.Name, programMessage.MessageID, content)
				if err != nil {
					return err
				}
				err = dtv.queries.UpdateProgram(ctx, p)
				if err != nil {
					return err
				}
				asp := NewAutoSearchProgram(p)
				for _, as := range autoSearchList {
					if as.IsMatchProgram(asp) {
						dtv.logger.Debug("program matched", zap.String("program.Name", p.Name), zap.String("as.Title", as.Title))
						err := dtv.sendAutoSearchMatchMessage(ctx, msg, p, service, as)
						if err != nil {
							dtv.logger.Warn("sendAutoSearchMatchMessage error", zap.Error(err))
							continue
						}
					}
				}
			}
		}
	}
	return nil
}

func (dtv *DTVUsecase) OnProgramsUpdated(ctx context.Context, serviceId uint) error {
	job, err := dtv.scheduler.Every("1m").LimitRunsTo(1).Do(func() {
		newCtx := context.Background()
		err := dtv.onProgramsUpdated(newCtx, serviceId)
		if err != nil {
			dtv.logger.Error("onProgramsUpdated error", zap.Error(err))
		}
		dtv.logger.Debug("onProgramsUpdated completed", zap.Uint("serviceId", serviceId))
	})
	if err != nil {
		dtv.logger.Error("scheduling error", zap.Error(err))
		return err
	}
	dtv.logger.Debug("scheduled onProgramsUpdated", zap.Uint("serviceId", serviceId), zap.Time("NextRun", job.NextRun()))
	return nil
}

func (dtv *DTVUsecase) sendAutoSearchMatchMessage(ctx context.Context, msg *discordgo.Message, p db.Program, service *db.Service, as *AutoSearch) error {
	url := discord.BuildMessageLinkURL(dtv.discord.Session().State.Guilds[0].ID, msg.ChannelID, msg.ID)
	content, err := template.GetAutoSearchMessage(p, *service, url)
	if err != nil {
		return err
	}
	content = width.Fold.String(content)
	notifyUsers, err := dtv.discord.GetMessageReactions(msg.ChannelID, msg.ID, discord.NotifyReactionEmoji)
	if err != nil {
		return err
	}
	notifyString := ""
	if len(notifyUsers) > 0 {
		for _, u := range notifyUsers {
			if u.ID != dtv.discord.Session().State.User.ID {
				notifyString += "<@" + u.ID + "> "
			}
		}
		notifyString += "\n"
	}
	content += notifyString
	err = dtv.discord.SendMessageToThread(as.ThreadID, content)
	if err != nil {
		return err
	}
	return nil
}
