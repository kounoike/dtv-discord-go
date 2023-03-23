package dtv

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/template"
	"github.com/pkg/errors"
)

func (dtv *DTVUsecase) OnProgramsUpdated(serviceId uint) {
	service, err := dtv.mirakc.GetService(serviceId)
	_ = service
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = dtv.discord.SendMessage("録画-情報", "動作ログ", fmt.Sprintf("programs updated: %s", service.Name))
	if err != nil {
		fmt.Println(err)
		return
	}
	programs, err := dtv.mirakc.ListPrograms(serviceId)
	if err != nil {
		fmt.Println(err)
		return
	}
	// if len(programs) > 0 {
	// 	fmt.Println(channel, len(programs), programs[0])
	// }
	ctx := context.Background()
	for _, p := range programs {
		if p.Name == "" {
			continue
		}
		program, err := dtv.queries.GetProgram(ctx, p.ID)
		if errors.Cause(err) == sql.ErrNoRows {
			msg, err := template.GetProgramMessage(p, *service)
			if err != nil {
				fmt.Println(err)
				return
			}
			msgID, err := dtv.discord.SendMessage("録画-番組情報", service.Name, msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			p.InsertDb(ctx, *dtv.queries)
			dtv.queries.InsertProgramMessage(ctx, db.InsertProgramMessageParams{MessageID: msgID, ProgramID: p.ID})
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
				p.UpdateDb(ctx, *dtv.queries)
			}
		}
	}
}
