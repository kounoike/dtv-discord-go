package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kounoike/dtv-discord-go/gpt"
	"github.com/mattn/go-shellwords"
	"go.uber.org/zap"
)

const TypeProgramTranscription = "program:transcription"

type ProgramTranscriptionPayload struct {
	ProgramId   int64  `json:"programId"`
	ContentPath string `json:"contentPath"`
	OutputPath  string `json:"outputPath"`
}

func NewProgramTranscriptionTask(programId int64, contentPath string, outputPath string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProgramTranscriptionPayload{ProgramId: programId, ContentPath: contentPath, OutputPath: outputPath})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProgramTranscription, payload, asynq.MaxRetry(10), asynq.Timeout(20*time.Hour), asynq.Retention(30*time.Minute)), nil
}

type ProgramTranscriber struct {
	logger              *zap.Logger
	recordedBasePath    string
	transcribedBasePath string
	gpt                 *gpt.GPTClient
}

func NewProgramTranscriber(logger *zap.Logger, gpt *gpt.GPTClient, encodeCommandTmpl *template.Template, recordedBasePath string, transcribedBasePath string) *ProgramTranscriber {
	return &ProgramTranscriber{
		logger:              logger,
		gpt:                 gpt,
		recordedBasePath:    recordedBasePath,
		transcribedBasePath: transcribedBasePath,
	}
}

func (e *ProgramTranscriber) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ProgramTranscriptionPayload
	err := json.Unmarshal(t.Payload(), &p)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	e.logger.Debug("Running ProcessTask", zap.Int64("programId", p.ProgramId), zap.String("contentPath", p.ContentPath), zap.String("outputPath", p.OutputPath))

	if p.ContentPath == "" || p.OutputPath == "" {
		e.logger.Error("empty ContentPath or OutputPath")
		return nil
	}

	tmpFile := fmt.Sprintf("/tmp/%d.m4a", p.ProgramId)
	commandLine := fmt.Sprintf(`ffmpeg -i "%s" -vn -ac 1 -ar 16000 -ab 32k "%s" -y`, path.Join(e.recordedBasePath, p.ContentPath), tmpFile)

	e.logger.Info("Running split audio command", zap.String("command", commandLine))

	args, err := shellwords.Parse(commandLine)
	if err != nil {
		return fmt.Errorf("split audio command shell parse error: %v: %w", err, asynq.SkipRetry)
	}

	var cmd *exec.Cmd
	switch len(args) {
	case 0:
		return fmt.Errorf("split audio command is empty %w", asynq.SkipRetry)
	case 1:
		cmd = exec.Command(args[0])
	default:
		cmd = exec.Command(args[0], args[1:]...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		e.logger.Error("split audio command execution error", zap.Error(err), zap.ByteString("output", out))
		return err
	}
	e.logger.Debug("split audio command succeeded")

	text, err := e.gpt.TranscribeText(ctx, tmpFile)
	os.Remove(tmpFile)
	if err != nil {
		e.logger.Error("TranscribeText error", zap.Error(err))
		return err
	}

	file, err := os.Create(path.Join(e.transcribedBasePath, p.OutputPath))
	if err != nil {
		e.logger.Error("File Create error", zap.Error(err))
		return err
	}

	_, err = file.WriteString(text)
	if err != nil {
		e.logger.Error("write transcribed text to file error", zap.Error(err))
		return err
	}
	err = file.Close()
	if err != nil {
		e.logger.Error("close error", zap.Error(err))
		return err
	}

	return nil
}
