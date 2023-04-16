package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const TypeProgramDeleteoriginal = "program:delete_original"

type ProgramDeleteOriginalPayload struct {
	ProgramId      int64    `json:"programId"`
	ContentPath    string   `json:"contentPath"`
	MonitorTaskIds []string `json:"monitorTaskIds"`
}

func NewProgramDeleteOriginalTask(programId int64, contentPath string, monitorTaskIds []string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProgramDeleteOriginalPayload{ProgramId: programId, ContentPath: contentPath, MonitorTaskIds: monitorTaskIds})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProgramDeleteoriginal, payload, asynq.MaxRetry(10), asynq.Timeout(20*time.Hour), asynq.Retention(30*time.Minute)), nil
}

type ProgramDeleter struct {
	logger           *zap.Logger
	inspector        *asynq.Inspector
	recordedBasePath string
}

func NewProgramDeleter(logger *zap.Logger, inspector *asynq.Inspector, recordedBasePath string) *ProgramDeleter {
	return &ProgramDeleter{
		logger:           logger,
		inspector:        inspector,
		recordedBasePath: recordedBasePath,
	}
}

func (e *ProgramDeleter) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ProgramDeleteOriginalPayload
	err := json.Unmarshal(t.Payload(), &p)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	e.logger.Debug("Running ProcessTask", zap.Int64("programId", p.ProgramId), zap.String("contentPath", p.ContentPath))

	if p.ContentPath == "" {
		return fmt.Errorf("empty ContentPath delete original command failed: %w", asynq.SkipRetry)
	}

	for _, monitorTaskId := range p.MonitorTaskIds {
		taskInfo, err := e.inspector.GetTaskInfo("default", monitorTaskId)
		if err != nil {
			e.logger.Error("get task failed", zap.Error(err))
			return fmt.Errorf("get task failed: %w", asynq.SkipRetry)
		}
		if taskInfo != nil && taskInfo.State == asynq.TaskStateArchived {
			e.logger.Error("task %s is failed", zap.String("taskType", taskInfo.Type))
			return fmt.Errorf("task %s is failed: %w", taskInfo.Type, asynq.SkipRetry)
		}
		if taskInfo != nil && taskInfo.State != asynq.TaskStateCompleted {
			return fmt.Errorf("task %s is running", taskInfo.Type)
		}
	}

	if err := os.Remove(filepath.Join(e.recordedBasePath, p.ContentPath)); err != nil {
		e.logger.Error("delete original file command failed", zap.Error(err))
		return fmt.Errorf("delete original file command failed: %w", asynq.SkipRetry)
	}

	e.logger.Debug("delete original file command succeeded")

	return nil
}
