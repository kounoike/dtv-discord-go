package gpt

import (
	"context"
	"encoding/json"

	"github.com/kounoike/dtv-discord-go/template"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

const systemPrompt = `Parse user message with this format:
{"title": title name of the program, "subtitle": subtitle name of the program including the part used to extract the number of episodes, empty if it does not exist, "episode": number of episodes extracted from Arabic or Chinese numerals in numeric type, 0 if not present}
Please return in strict JSON format. Never include non-JSON content in the output, such as commentary. Line breaks are not required.`

type GPTClient struct {
	enabled bool
	token   string
	logger  *zap.Logger
}

func NewGPTClient(enabled bool, token string, logger *zap.Logger) *GPTClient {
	return &GPTClient{
		enabled: enabled,
		token:   token,
		logger:  logger,
	}
}

func (c *GPTClient) ParseTitle(ctx context.Context, title string, pathTemplateData *template.PathTemplateData) error {
	pathTemplateData.Title = title
	if c.enabled {
		client := openai.NewClient(c.token)
		req := openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: title,
				},
			},
		}
		resp, err := client.CreateChatCompletion(ctx, req)
		if err != nil {
			return err
		}
		c.logger.Debug("Success ChatComplettion Request #1", zap.String("response", resp.Choices[0].Message.Content))

		err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), pathTemplateData)
		if err != nil {
			resp, err := client.CreateChatCompletion(ctx, req)
			if err != nil {
				return err
			}
			c.logger.Debug("Success ChatComplettion Request #2", zap.String("response", resp.Choices[0].Message.Content))
			err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), pathTemplateData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *GPTClient) TranscribeText(ctx context.Context, audioFilePath string) (string, error) {
	if !c.enabled {
		return "", nil
	}
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		Language: "ja",
		FilePath: audioFilePath,
	}
	client := openai.NewClient(c.token)
	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Text, nil
}
