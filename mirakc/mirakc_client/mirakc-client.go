package mirakc_client

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_model"
	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

type MirakcClient struct {
	host   string
	port   uint
	logger *zap.Logger
}

func NewMirakcClient(host string, port uint, logger *zap.Logger) *MirakcClient {
	return &MirakcClient{host: host, port: port, logger: logger}
}

func (m *MirakcClient) GetVersion() (*mirakc_model.Version, error) {
	url := fmt.Sprintf("http://%s:%d/api/version", m.host, m.port)
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP Error status code: %d", resp.StatusCode())
	}
	var version mirakc_model.Version
	if err = json.Unmarshal(resp.Body(), &version); err != nil {
		return nil, err
	}
	return &version, nil
}

func (m *MirakcClient) ListServices() ([]db.Service, error) {
	url := fmt.Sprintf("http://%s:%d/api/services", m.host, m.port)
	client := resty.New()
	resp, err := client.R().
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP Error status code: %d", resp.StatusCode())
	}
	var services []db.Service
	if err = json.Unmarshal(resp.Body(), &services); err != nil {
		return nil, err
	}
	return services, nil
}

func (m *MirakcClient) GetService(serviceId uint) (*db.Service, error) {
	url := fmt.Sprintf("http://%s:%d/api/services/%d", m.host, m.port, serviceId)
	client := resty.New()
	resp, err := client.R().
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP Error status code: %d", resp.StatusCode())
	}

	var service db.Service

	if err = json.Unmarshal(resp.Body(), &service); err != nil {
		return nil, err
	}
	return &service, nil
}

func (m *MirakcClient) ListPrograms(serviceId uint) ([]db.Program, error) {
	url := fmt.Sprintf("http://%s:%d/api/services/%d/programs", m.host, m.port, serviceId)
	client := resty.New()
	resp, err := client.R().
		// EnableTrace().
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP Error status code: %d", resp.StatusCode())
	}

	any := jsoniter.Get(resp.Body())
	programs := make([]db.Program, 0, any.Size())
	for i := 0; i < any.Size(); i++ {
		jsonStr := any.Get(i).ToString()
		var program db.Program
		if err := json.Unmarshal([]byte(jsonStr), &program); err != nil {
			continue
		}
		programs = append(programs, program)
	}

	return programs, nil
}

type scheduleOptions struct {
	ContentPath string `json:"contentPath"`
}

type scheduleData struct {
	ProgramID int64           `json:"programId"`
	Options   scheduleOptions `json:"options"`
	Tags      []string        `json:"tags"`
}

func (m *MirakcClient) AddRecordingSchedule(programID int64, contentPath string) error {
	url := fmt.Sprintf("http://%s:%d/api/recording/schedules", m.host, m.port)
	data := scheduleData{
		ProgramID: programID,
		Options: scheduleOptions{
			ContentPath: contentPath,
		},
		Tags: []string{"dtv-discord-bot"},
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// postOption := fmt.Sprintf(`{"programId": %d, "options": {"contentPath": "%d.m2ts"}, "tags": ["manual"]}`, programID, programID)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(dataJson).
		Post(url)
	if err != nil {
		return err
	}
	m.logger.Info("録画予約完了", zap.Int("StatusCode", resp.StatusCode()), zap.String("contentPath", contentPath))
	if resp.StatusCode() == 201 {
		return nil
	}
	return fmt.Errorf("post request:%s with body:%s status code:%d", url, string(dataJson), resp.StatusCode())
}

func (m *MirakcClient) DeleteRecordingSchedule(programID int64) error {
	url := fmt.Sprintf("http://%s:%d/api/recording/schedules", m.host, m.port)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Delete(url)
	if err != nil {
		return err
	}
	if resp.StatusCode()/100 == 2 {
		return nil
	}
	return fmt.Errorf("delete request:%s status code:%d", url, resp.StatusCode())
}

func (m *MirakcClient) GetRecordingScheduleContentPath(programID int64) (string, error) {
	url := fmt.Sprintf("http://%s:%d/api/recording/schedules/%d", m.host, m.port, programID)
	client := resty.New()
	resp, err := client.R().
		Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode()/100 == 2 {
		var data scheduleData
		err := json.Unmarshal(resp.Body(), &data)
		if err != nil {
			return "", err
		}
		return data.Options.ContentPath, nil
	}
	return "", fmt.Errorf("get request:%s status code:%d", url, resp.StatusCode())
}
