package meili

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
	"golang.org/x/text/width"
)

type MeiliSearchClient struct {
	logger              *zap.Logger
	client              meilisearch.ClientInterface
	transcribedBasePath string
}

func NewMeiliSearchClient(logger *zap.Logger, host string, port int, transcribedBasePath string) *MeiliSearchClient {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: fmt.Sprintf("http://%s:%d", host, port),
	})
	return &MeiliSearchClient{logger, client, transcribedBasePath}
}

func (m *MeiliSearchClient) Index(name string) *meilisearch.Index {
	return m.client.Index(name)
}

func (m *MeiliSearchClient) UpdatePrograms(programs []db.ListProgramWithMessageAndServiceNameRow, guildID string) error {
	index := m.Index("program")
	_, err := index.UpdateFilterableAttributes(&[]string{"チャンネル名", "ジャンル"})
	if err != nil {
		return err
	}
	for _, program := range programs {
		if program.Name == "" {
			continue
		}
		discordMessageUrl := fmt.Sprintf("discord://discord.com/channels/%s/%s/%s", guildID, program.ChannelID, program.MessageID)
		webMessageUrl := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, program.ChannelID, program.MessageID)
		document := map[string]interface{}{
			"id":                program.ProgramID,
			"タイトル":              width.Fold.String(program.Name),
			"番組説明":              width.Fold.String(program.Description),
			"ジャンル":              width.Fold.String(program.Genre),
			"番組詳細":              width.Fold.String(getProgramExtendedFromJson(program.Json)),
			"チャンネル名":            width.Fold.String(program.ServiceName),
			"WebMessageUrl":     webMessageUrl,
			"DiscordMessageUrl": discordMessageUrl,
			"StartAt":           program.StartAt,
			"Duration":          program.Duration,
		}
		_, err := index.UpdateDocuments([]map[string]interface{}{document})
		if err != nil {
			m.logger.Warn("failed to update documents", zap.Error(err), zap.Any("document", document))
		}
	}

	return nil
}

func (m *MeiliSearchClient) DeleteProgramIndex() error {
	_, err := m.client.DeleteIndex("program")
	return err
}

func (m *MeiliSearchClient) DeleteRecordedFileIndex() error {
	_, err := m.client.DeleteIndex("recorded_file")
	return err
}

func (m *MeiliSearchClient) UpdateRecordedFiles(rows []db.ListRecordedFilesRow) error {
	index := m.Index("recorded_file")
	_, err := index.UpdateFilterableAttributes(&[]string{"チャンネル名", "ジャンル"})
	if err != nil {
		return err
	}

	for _, row := range rows {
		document := map[string]interface{}{
			"id":       row.ProgramID,
			"タイトル":     width.Fold.String(row.Name),
			"番組説明":     width.Fold.String(row.Description),
			"ジャンル":     width.Fold.String(row.Genre),
			"番組詳細":     width.Fold.String(getProgramExtendedFromJson(row.Json)),
			"チャンネル名":   width.Fold.String(row.ServiceName),
			"StartAt":  row.StartAt,
			"Duration": row.Duration,
		}
		if row.Mp4Path.Valid {
			document["mp4"] = row.Mp4Path.String
		}
		if row.M2tsPath.Valid {
			document["m2ts"] = row.M2tsPath.String
		}
		if row.Aribb24TxtPath.Valid {
			bytes, err := ioutil.ReadFile(filepath.Join(m.transcribedBasePath, row.Aribb24TxtPath.String))
			if err != nil {
				m.logger.Warn("failed to read aribb24 txt", zap.Error(err), zap.String("path", row.Aribb24TxtPath.String))
			} else {
				document["ARIB字幕"] = string(bytes)
			}
		}
		if row.TranscribedTxtPath.Valid {
			bytes, err := ioutil.ReadFile(filepath.Join(m.transcribedBasePath, row.TranscribedTxtPath.String))
			if err != nil {
				m.logger.Warn("failed to read transcribed txt", zap.Error(err), zap.String("path", row.TranscribedTxtPath.String))
			} else {
				document["文字起こし"] = string(bytes)
			}
		}
		_, err := index.UpdateDocuments([]map[string]interface{}{document})
		if err != nil {
			m.logger.Warn("failed to update documents", zap.Error(err), zap.Any("document", document))
		}
	}

	return nil
}

func (m *MeiliSearchClient) UpdateRecordedFile(row db.GetRecordedFilesRow) error {
	document := map[string]interface{}{
		"id":       row.ProgramID,
		"タイトル":     width.Fold.String(row.Name),
		"番組説明":     width.Fold.String(row.Description),
		"ジャンル":     width.Fold.String(row.Genre),
		"番組詳細":     width.Fold.String(getProgramExtendedFromJson(row.Json)),
		"チャンネル名":   width.Fold.String(row.ServiceName),
		"StartAt":  row.StartAt,
		"Duration": row.Duration,
	}
	if row.Mp4Path.Valid {
		document["mp4"] = row.Mp4Path.String
	}
	if row.M2tsPath.Valid {
		document["m2ts"] = row.M2tsPath.String
	}
	if row.Aribb24TxtPath.Valid {
		bytes, err := ioutil.ReadFile(filepath.Join(m.transcribedBasePath, row.Aribb24TxtPath.String))
		if err != nil {
			m.logger.Warn("failed to read aribb24 txt", zap.Error(err), zap.String("path", row.Aribb24TxtPath.String))
		} else {
			document["ARIB字幕"] = string(bytes)
		}
	}
	if row.TranscribedTxtPath.Valid {
		bytes, err := ioutil.ReadFile(filepath.Join(m.transcribedBasePath, row.TranscribedTxtPath.String))
		if err != nil {
			m.logger.Warn("failed to read transcribed txt", zap.Error(err), zap.String("path", row.TranscribedTxtPath.String))
		} else {
			document["文字起こし"] = string(bytes)
		}
	}
	index := m.Index("recorded_file")
	_, err := index.UpdateDocuments(document)
	if err != nil {
		return err
	}
	_, err = index.UpdateFilterableAttributes(&[]string{"チャンネル名", "ジャンル"})
	if err != nil {
		return err
	}

	return nil
}

func getProgramExtendedFromJson(json string) string {
	b := []byte(json)
	any := jsoniter.Get(b, "extended")
	str := ""
	for _, key := range any.Keys() {
		str += fmt.Sprintf("%s: %s\n", key, any.Get(key).ToString())
	}
	return str
}
