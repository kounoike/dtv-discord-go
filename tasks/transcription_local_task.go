package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mattn/go-shellwords"
	"go.uber.org/zap"
)

const TypeProgramTranscriptionLocal = "program:transcription:local"

type ProgramTranscriptionLocalPayload struct {
	ProgramId   int64  `json:"programId"`
	ContentPath string `json:"contentPath"`
	EncodedPath string `json:"encodedPath"`
	OutputPath  string `json:"outputPath"`
}

func NewProgramTranscriptionLocalTask(programId int64, contentPath string, encodedPath string, outputPath string, queueName string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProgramTranscriptionLocalPayload{ProgramId: programId, ContentPath: contentPath, EncodedPath: encodedPath, OutputPath: outputPath})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProgramTranscriptionLocal, payload, asynq.MaxRetry(10), asynq.Timeout(20*time.Hour), asynq.Retention(30*time.Minute), asynq.Queue(queueName)), nil
}

type ProgramTranscriberLocal struct {
	logger              *zap.Logger
	recordedBasePath    string
	encodedBasePath     string
	transcribedBasePath string
	runWhisperScript    string
	whisperModel        string
}

func NewProgramTranscriberLocal(logger *zap.Logger, encodeCommandTmpl *template.Template, recordedBasePath string, encodedBasePath string, transcribedBasePath string, runWhisperScript string, whisperModel string) *ProgramTranscriberLocal {
	return &ProgramTranscriberLocal{
		logger:              logger,
		recordedBasePath:    recordedBasePath,
		encodedBasePath:     encodedBasePath,
		transcribedBasePath: transcribedBasePath,
		runWhisperScript:    runWhisperScript,
		whisperModel:        whisperModel,
	}
}

func (e *ProgramTranscriberLocal) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ProgramTranscriptionLocalPayload
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

	inputFile := path.Join(e.recordedBasePath, p.ContentPath)
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		inputFile = path.Join(e.encodedBasePath, p.EncodedPath)
	}

	tmpFile := fmt.Sprintf("/tmp/%d.wav", p.ProgramId)
	commandLine := fmt.Sprintf(`ffmpeg -hide_banner -i "%s" -vn "%s" -y`, inputFile, tmpFile)

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

	whisperCommandLine := fmt.Sprintf(`python "%s" "%s" "%s" "%s"`, e.runWhisperScript, e.whisperModel, tmpFile, path.Join(e.transcribedBasePath, p.OutputPath))

	e.logger.Info("Running whisper script", zap.String("command", whisperCommandLine))

	whisperArgs, err := shellwords.Parse(whisperCommandLine)
	if err != nil {
		return fmt.Errorf("whisper command shell parse error: %v: %w", err, asynq.SkipRetry)
	}

	var whisperCmd *exec.Cmd
	switch len(whisperArgs) {
	case 0:
		return fmt.Errorf("whisper command is empty %w", asynq.SkipRetry)
	case 1:
		whisperCmd = exec.Command(whisperArgs[0])
	default:
		whisperCmd = exec.Command(whisperArgs[0], whisperArgs[1:]...)
	}
	whisperOut, err := whisperCmd.CombinedOutput()
	if err != nil {
		e.logger.Error("whisper command execution error", zap.Error(err), zap.ByteString("output", whisperOut))
		return err
	}
	e.logger.Debug("whisper command succeeded")

	return nil
}
