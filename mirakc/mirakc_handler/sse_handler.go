package mirakc_handler

import (
	"encoding/json"
	"fmt"

	"github.com/kounoike/dtv-discord-go/dtv"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_model"
	"github.com/r3labs/sse/v2"
)

type SSEHandler struct {
	dtv dtv.DTVUsecase
	sse *sse.Client
}

func NewSSEHandler(dtv dtv.DTVUsecase, host string, port uint) *SSEHandler {
	sseClient := sse.NewClient(fmt.Sprintf("http://%s:%d/events", host, port))

	return &SSEHandler{
		dtv: dtv,
		sse: sseClient,
	}
}

func (h *SSEHandler) onProgramsUpdated(serviceId uint) {
	h.dtv.OnProgramsUpdated(serviceId)
}

func (h *SSEHandler) Subscribe() {
	h.sse.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		eventName := string(msg.Event)
		fmt.Printf("%s: %s\n", eventName, string(msg.Data))
		if eventName == "epg.programs-updated" {
			var data mirakc_model.ProgramsUpdatedEventData
			err := json.Unmarshal(msg.Data, &data)
			if err != nil {
				fmt.Println(err)
				return
			}
			h.onProgramsUpdated(data.ServiceId)
		}
	})
}
