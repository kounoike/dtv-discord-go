package mirakc_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kounoike/dtv-discord-go/dtv"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_model"
	"github.com/r3labs/sse/v2"
	"go.uber.org/zap"
)

type SSEHandler struct {
	dtv    dtv.DTVUsecase
	sse    *sse.Client
	logger *zap.Logger
}

func NewSSEHandler(dtv dtv.DTVUsecase, host string, port uint, logger *zap.Logger) *SSEHandler {
	sseClient := sse.NewClient(fmt.Sprintf("http://%s:%d/events", host, port))

	return &SSEHandler{
		dtv:    dtv,
		sse:    sseClient,
		logger: logger,
	}
}

func (h *SSEHandler) onProgramsUpdated(serviceId uint) {
	ctx := context.Background()
	err := h.dtv.OnProgramsUpdated(ctx, serviceId)
	if err != nil {
		h.logger.Error("OnProgramsUpdate error", zap.Error(err))
	}
}

func (h *SSEHandler) Subscribe() {
	h.sse.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		eventName := string(msg.Event)
		h.logger.Debug("sse event received", zap.String("eventName", eventName), zap.String("Data", string(msg.Data)))
		if eventName == "epg.programs-updated" {
			var data mirakc_model.ProgramsUpdatedEventData
			err := json.Unmarshal(msg.Data, &data)
			if err != nil {
				h.logger.Error("json Unmarshal error", zap.Error(err))
				return
			}
			h.onProgramsUpdated(data.ServiceId)
		}
		h.logger.Debug("sse event processed successfully")
	})
}
