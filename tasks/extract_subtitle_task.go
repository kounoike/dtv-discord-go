package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mattn/go-shellwords"
	"go.uber.org/zap"
)

const TypeProgramExtractSubtitle = "program:extract_subtitle"

type ProgramExtractSubtilePayload struct {
	ProgramId   int64  `json:"programId"`
	ContentPath string `json:"contentPath"`
	OutputPath  string `json:"outputPath"`
}

func NewProgramExtractSubtileTask(programId int64, contentPath string, outputPath string, queueName string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProgramExtractSubtilePayload{ProgramId: programId, ContentPath: contentPath, OutputPath: outputPath})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProgramExtractSubtitle, payload, asynq.MaxRetry(10), asynq.Timeout(20*time.Hour), asynq.Retention(3*24*time.Hour), asynq.Queue(queueName)), nil
}

type ProgramExtractor struct {
	logger              *zap.Logger
	recordedBasePath    string
	transcribedBasePath string
}

func NewProgramExtractor(logger *zap.Logger, recordedBasePath string, transcribedBasePath string) *ProgramExtractor {
	return &ProgramExtractor{
		logger:              logger,
		recordedBasePath:    recordedBasePath,
		transcribedBasePath: transcribedBasePath,
	}
}

func (e *ProgramExtractor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ProgramEncodePayload
	err := json.Unmarshal(t.Payload(), &p)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	e.logger.Debug("Running ProcessTask", zap.Int64("programId", p.ProgramId), zap.String("contentPath", p.ContentPath), zap.String("outputPath", p.OutputPath))

	if p.ContentPath == "" || p.OutputPath == "" {
		e.logger.Error("empty ContentPath or OutputPath")
		return nil
	}

	err = os.MkdirAll(filepath.Dir(filepath.Join(e.transcribedBasePath, p.OutputPath)), 0777)
	if err != nil {
		return err
	}

	commandLine := fmt.Sprintf(`ffmpeg -i "%s" -an -vn -c:s text -f rawvideo "%s"`, filepath.Join(e.recordedBasePath, p.ContentPath), filepath.Join(e.transcribedBasePath, p.OutputPath))

	e.logger.Info("Running extract command", zap.String("command", commandLine))

	args, err := shellwords.Parse(commandLine)
	if err != nil {
		return fmt.Errorf("extract command shell parse error: %v: %w", err, asynq.SkipRetry)
	}

	var cmd *exec.Cmd
	switch len(args) {
	case 0:
		return fmt.Errorf("extract command is empty %w", asynq.SkipRetry)
	case 1:
		cmd = exec.Command(args[0])
	default:
		cmd = exec.Command(args[0], args[1:]...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		e.logger.Error("extract command execution error", zap.Error(err), zap.ByteString("output", out))
		return err
	}
	e.logger.Debug("extract command succeeded")

	return nil
}
