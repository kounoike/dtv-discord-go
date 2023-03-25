package dtv

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/template"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/width"
)

func (dtv *DTVUsecase) OnProgramsUpdated(ctx context.Context, serviceId uint) error {
	service, err := dtv.mirakc.GetService(serviceId)
	_ = service
	if err != nil {
		return err
	}
	_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.LogChannel, fmt.Sprintf("programs updated: %s", service.Name))
	if err != nil {
		return err
	}
	programs, err := dtv.mirakc.ListPrograms(serviceId)
	if err != nil {
		return err
	}

	autoSearchList, err := dtv.ListAutoSearchForServiceName(service.Name)
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
			err = dtv.queries.InsertProgram(ctx, p)
			if err != nil {
				return err
			}
			err = dtv.queries.InsertProgramMessage(ctx, db.InsertProgramMessageParams{MessageID: msg.ID, ProgramID: p.ID})
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
				dtv.logger.Debug("matching", zap.String("p.Name", p.Name), zap.String("asp.Title", asp.Title), zap.String("as.Title", as.Title), zap.Bool("isMatch", as.IsMatchProgram(asp)))
				if as.IsMatchProgram(asp) {
					dtv.logger.Debug("program matched", zap.String("program.Name", p.Name), zap.String("as.Title", as.Title))
					url := discord.BuildMessageLinkURL(dtv.discord.Session().State.Guilds[0].ID, msg.ChannelID, msg.ID)
					content, err := template.GetAutoSearchMessage(p, *service, url)
					if err != nil {
						return err
					}
					content = width.Fold.String(content)
					notifyString := ""
					recorderString := ""
					if len(as.NotifyUsers) > 0 {
						for _, u := range as.NotifyUsers {
							notifyString += "<@" + u.ID + "> "
						}
						notifyString += "\n"
					}
					if len(as.RecordingUsers) > 0 {
						for _, u := range as.RecordingUsers {
							recorderString += "<@" + u.ID + "> "
						}
						recorderString += "録画予約しました\n"
					}
					content += notifyString + recorderString
					err = dtv.discord.SendMessageToThread(as.ThreadID, content)
					if err != nil {
						return err
					}
					if len(as.RecordingUsers) > 0 {
						dtv.discord.MessageReactionAdd(msg.ChannelID, msg.ID, discord.RecordingReactionEmoji)
						dtv.checkRecordingForMessage(ctx, msg.ChannelID, msg.ID)
					}
				}
			}
		} else if err != nil {
			return err
		} else {
			pJson, err := p.Json.MarshalJSON()
			if err != nil {
				continue
			}
			programJson, err := program.Json.MarshalJSON()
			if err != nil {
				continue
			}
			if !bytes.Equal(pJson, programJson) {
				dtv.queries.UpdateProgram(ctx, p)
			}
		}
	}
	return nil
}
