package sse_handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/template"
	"github.com/kounoike/dtv-discord-go/tv"
	"github.com/pingcap/errors"
	"github.com/r3labs/sse/v2"
)

type SSEHandler struct {
	ctx     context.Context
	mirakc  tv.MirakcClient
	discord discord.DiscordClient
	sse     sse.Client
	queries db.Queries
}

func NewSSEHandler(ctx context.Context, mirakc tv.MirakcClient, discord discord.DiscordClient, sse sse.Client, queries db.Queries) *SSEHandler {
	return &SSEHandler{
		ctx:     ctx,
		discord: discord,
		mirakc:  mirakc,
		sse:     sse,
		queries: queries,
	}
}

func (h *SSEHandler) onProgramsUpdated(serviceId uint) {
	service, err := h.mirakc.GetService(serviceId)
	_ = service
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = h.discord.SendMessage("録画-情報", "動作ログ", fmt.Sprintf("programs updated: %s", service.Name))
	if err != nil {
		fmt.Println(err)
		return
	}
	programs, err := h.mirakc.ListPrograms(serviceId)
	if err != nil {
		fmt.Println(err)
		return
	}
	// if len(programs) > 0 {
	// 	fmt.Println(channel, len(programs), programs[0])
	// }
	for _, p := range programs {
		if p.Name == "" {
			continue
		}
		program, err := h.queries.GetProgram(h.ctx, p.ID)
		if errors.Cause(err) == sql.ErrNoRows {
			msg, err := template.GetProgramMessage(p, *service)
			if err != nil {
				fmt.Println(err)
				return
			}
			msgID, err := h.discord.SendMessage("録画-番組情報", service.Name, msg)
			p.InsertDb(h.ctx, h.queries)
			h.queries.InsertProgramMessage(h.ctx, db.InsertProgramMessageParams{MessageID: msgID, ProgramID: p.ID})
		} else {
			pJson, err := p.Json.MarshalJSON()
			if err != nil {
				continue
			}
			programJson, err := program.Json.MarshalJSON()
			if err != nil {
				continue
			}
			if bytes.Compare(pJson, programJson) != 0 {
				p.UpdateDb(h.ctx, h.queries)
			}
		}
	}
}

func (h *SSEHandler) Subscribe() {
	h.sse.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		eventName := string(msg.Event)
		fmt.Printf("%s: %s\n", eventName, string(msg.Data))
		if eventName == "epg.programs-updated" {
			var data tv.ProgramsUpdatedEventData
			err := json.Unmarshal(msg.Data, &data)
			if err != nil {
				fmt.Println(err)
				return
			}
			h.onProgramsUpdated(data.ServiceId)
		}
	})
}
