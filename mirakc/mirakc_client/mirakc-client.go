package mirakc_client

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_model"

	"github.com/go-resty/resty/v2"
)

type MirakcClient struct {
	host string
	port uint
}

func NewMirakcClient(host string, port uint) *MirakcClient {
	return &MirakcClient{host: host, port: port}
}

func (m *MirakcClient) GetService(serviceId uint) (*mirakc_model.Service, error) {
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

	var service mirakc_model.Service

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

	// // Explore response object
	// fmt.Println("Response Info:")
	// fmt.Println("  Error      :", err)
	// fmt.Println("  Status Code:", resp.StatusCode())
	// fmt.Println("  Status     :", resp.Status())
	// fmt.Println("  Proto      :", resp.Proto())
	// fmt.Println("  Time       :", resp.Time())
	// fmt.Println("  Received At:", resp.ReceivedAt())
	// // fmt.Println("  Body       :\n", resp)
	// fmt.Println()

	// // Explore trace info
	// fmt.Println("Request Trace Info:")
	// ti := resp.Request.TraceInfo()
	// fmt.Println("  DNSLookup     :", ti.DNSLookup)
	// fmt.Println("  ConnTime      :", ti.ConnTime)
	// fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	// fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	// fmt.Println("  ServerTime    :", ti.ServerTime)
	// fmt.Println("  ResponseTime  :", ti.ResponseTime)
	// fmt.Println("  TotalTime     :", ti.TotalTime)
	// fmt.Println("  IsConnReused  :", ti.IsConnReused)
	// fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	// fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	// fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	// fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())

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

	// fmt.Printf("%+v\n", programs[0])
	// fmt.Println("Start:", programs[0].GetStartTime().Local().Format(time.RFC1123))
	// fmt.Println("End:", programs[0].GetEndTime().Local().Format(time.RFC1123))

	return programs, nil
}

func (m *MirakcClient) AddRecordingSchedule(programID int64) error {
	url := fmt.Sprintf("http://%s:%d/api/recording/schedules", m.host, m.port)
	postOption := fmt.Sprintf(`{"programId": %d, "options": {"contentPath": "%d.m2ts"}, "tags": ["manual"]}`, programID, programID)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(postOption).
		Post(url)
	if err != nil {
		return err
	}
	fmt.Println("録画予約: StatusCode:", resp.StatusCode())
	if resp.StatusCode() == 201 {
		return nil
	}
	return fmt.Errorf("post request:%s status code:%d", url, resp.StatusCode())
}
