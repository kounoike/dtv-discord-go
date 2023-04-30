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

const (
	programIndexName      = "program"
	recordedFileIndexName = "recorded_file"
)

func NewMeiliSearchClient(logger *zap.Logger, host string, port int, transcribedBasePath string) *MeiliSearchClient {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: fmt.Sprintf("http://%s:%d", host, port),
	})
	return &MeiliSearchClient{logger, client, transcribedBasePath}
}

func (m *MeiliSearchClient) Init() error {
	programIndex := m.Index(programIndexName)
	recordedFileIndex := m.Index(recordedFileIndexName)
	if _, err := programIndex.UpdateSearchableAttributes(&[]string{"タイトル", "番組説明", "ジャンル", "番組詳細", "チャンネル名"}); err != nil {
		return err
	}
	if _, err := recordedFileIndex.UpdateSearchableAttributes(&[]string{"タイトル", "番組説明", "ジャンル", "番組詳細", "チャンネル名", "ARIB字幕", "文字起こし"}); err != nil {
		return err
	}
	if _, err := programIndex.UpdateFilterableAttributes(&[]string{"チャンネル名", "ジャンル"}); err != nil {
		return err
	}
	if _, err := recordedFileIndex.UpdateFilterableAttributes(&[]string{"チャンネル名", "ジャンル"}); err != nil {
		return err
	}
	if _, err := programIndex.UpdateIndex("id"); err != nil {
		return err
	}
	if _, err := recordedFileIndex.UpdateIndex("id"); err != nil {
		return err
	}
	return nil
}

func (m *MeiliSearchClient) Index(name string) *meilisearch.Index {
	return m.client.Index(name)
}

func (m *MeiliSearchClient) UpdatePrograms(programs []db.ListProgramWithMessageAndServiceNameRow, guildID string) error {
	index := m.Index(programIndexName)
	documents := make([]map[string]interface{}, 0, len(programs))
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
		documents = append(documents, document)
	}

	if _, err := index.UpdateDocuments(documents); err != nil {
		m.logger.Warn("failed to update documents", zap.Error(err))
	}

	return nil
}

func (m *MeiliSearchClient) DeleteProgramIndex() error {
	if _, err := m.client.DeleteIndex(programIndexName); err != nil {
		return err
	}
	return m.Init()
}

func (m *MeiliSearchClient) DeleteRecordedFileIndex() error {
	if _, err := m.client.DeleteIndex(recordedFileIndexName); err != nil {
		return err
	}
	return m.Init()
}

func (m *MeiliSearchClient) UpdateRecordedFiles(rows []db.ListRecordedFilesRow) error {
	index := m.Index(recordedFileIndexName)

	documents := make([]map[string]interface{}, 0, len(rows))
	for idx, row := range rows {
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
		documents = append(documents, document)
		if idx%100 == 0 {
			m.logger.Info(fmt.Sprintf("%d/%d 番組 準備完了...", idx, len(rows)))
		}
	}

	m.logger.Info(fmt.Sprintf("%d/%d 番組 準備完了...", len(rows), len(rows)))

	_, err := index.UpdateDocuments(documents)
	if err != nil {
		m.logger.Warn("failed to update documents", zap.Error(err), zap.Any("documents", documents))
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
	index := m.Index(recordedFileIndexName)
	_, err := index.UpdateDocuments([]map[string]interface{}{document})
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
