package mirakc_client

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/kounoike/dtv-discord-go/db"
	"golang.org/x/exp/slog"

	"github.com/go-resty/resty/v2"
)

type MirakcClient struct {
	host string
	port uint
}

func NewMirakcClient(host string, port uint) *MirakcClient {
	return &MirakcClient{host: host, port: port}
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
		slog.Debug(string(resp.Body()))
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
		Tags: []string{"manual"},
	}
	// postOption := fmt.Sprintf(`{"programId": %d, "options": {"contentPath": "%d.m2ts"}, "tags": ["manual"]}`, programID, programID)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(url)
	if err != nil {
		return err
	}
	slog.Info("録画予約完了", "StatusCode", resp.StatusCode())
	if resp.StatusCode() == 201 {
		return nil
	}
	return fmt.Errorf("post request:%s status code:%d", url, resp.StatusCode())
}
