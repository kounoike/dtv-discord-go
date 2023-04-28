package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alessio/shellescape"
	"github.com/hibiken/asynq"
	"github.com/mattn/go-shellwords"
	"go.uber.org/zap"
)

const TypeProgramEncode = "program:encode"

type ProgramEncodePayload struct {
	ProgramId   int64  `json:"programId"`
	ContentPath string `json:"contentPath"`
	OutputPath  string `json:"outputPath"`
}

type commandTemplateData struct {
	InputPath  string
	OutputPath string
}

func NewProgramEncodeTask(programId int64, contentPath string, outputPath string, queueName string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProgramEncodePayload{ProgramId: programId, ContentPath: contentPath, OutputPath: outputPath})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProgramEncode, payload, asynq.MaxRetry(10), asynq.Timeout(20*time.Hour), asynq.Retention(3*24*time.Hour), asynq.Queue(queueName)), nil
}

type ProgramEncoder struct {
	logger                *zap.Logger
	encodeCommandTemplate *template.Template
	recordedBasePath      string
	encodedBasePath       string
	deleteOriginalFile    bool
}

func NewProgramEncoder(logger *zap.Logger, encodeCommandTmpl *template.Template, recordedBasePath string, encodedBasePath string, deleteOriginalFile bool) *ProgramEncoder {
	return &ProgramEncoder{
		logger:                logger,
		encodeCommandTemplate: encodeCommandTmpl,
		recordedBasePath:      recordedBasePath,
		encodedBasePath:       encodedBasePath,
		deleteOriginalFile:    deleteOriginalFile,
	}
}

func (e *ProgramEncoder) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ProgramEncodePayload
	err := json.Unmarshal(t.Payload(), &p)
	if err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	e.logger.Debug("Running ProcessTask", zap.Int64("programId", p.ProgramId), zap.String("contentPath", p.ContentPath), zap.String("outputPath", p.OutputPath))

	var buf bytes.Buffer

	if p.ContentPath == "" || p.OutputPath == "" {
		e.logger.Error("empty ContentPath or OutputPath")
		return nil
	}

	err = os.MkdirAll(filepath.Dir(filepath.Join(e.encodedBasePath, p.OutputPath)), 0777)
	if err != nil {
		return err
	}

	err = e.encodeCommandTemplate.Execute(&buf, commandTemplateData{
		InputPath:  shellescape.Quote(filepath.Join(e.recordedBasePath, p.ContentPath)),
		OutputPath: shellescape.Quote(filepath.Join(e.encodedBasePath, p.OutputPath)),
	})
	if err != nil {
		return fmt.Errorf("encode command template error: %v: %w", err, asynq.SkipRetry)
	}
	commandLine := buf.String()

	e.logger.Info("Running encode command", zap.String("command", commandLine))

	args, err := shellwords.Parse(commandLine)
	if err != nil {
		return fmt.Errorf("encode command shell parse error: %v: %w", err, asynq.SkipRetry)
	}

	var cmd *exec.Cmd
	switch len(args) {
	case 0:
		return fmt.Errorf("encode command is empty %w", asynq.SkipRetry)
	case 1:
		cmd = exec.Command(args[0])
	default:
		cmd = exec.Command(args[0], args[1:]...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		e.logger.Error("encode command execution error", zap.Error(err), zap.ByteString("output", out))
		t.ResultWriter().Write([]byte(out))
		return err
	}
	e.logger.Debug("encode command succeeded")

	return nil
}
