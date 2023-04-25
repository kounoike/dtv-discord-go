// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"database/sql"
	"time"
)

type ComponentVersion struct {
	ID        int32     `json:"id"`
	Component string    `json:"component"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DiscordServer struct {
	ID        int32     `json:"id"`
	ServerID  string    `json:"serverID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EncodeTask struct {
	ID        int32     `json:"id"`
	TaskID    string    `json:"taskID"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type IndexInvalid struct {
	ID        int32     `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Program struct {
	ID          int64     `json:"id"`
	Json        string    `json:"-"`
	EventID     int32     `json:"eventId"`
	ServiceID   int32     `json:"serviceId"`
	NetworkID   int32     `json:"networkId"`
	StartAt     int64     `json:"startAt"`
	Duration    int32     `json:"duration"`
	IsFree      bool      `json:"isFree"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Genre       string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type ProgramMessage struct {
	ID        int32     `json:"id"`
	ChannelID string    `json:"channelID"`
	MessageID string    `json:"messageID"`
	ProgramID int64     `json:"programID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProgramRecording struct {
	ID          int32     `json:"id"`
	ProgramID   int64     `json:"programID"`
	ContentPath string    `json:"contentPath"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProgramService struct {
	ID        int32     `json:"id"`
	ProgramID int64     `json:"programID"`
	ServiceID int64     `json:"serviceID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RecordedFile struct {
	ID                 int32          `json:"id"`
	ProgramID          int64          `json:"programID"`
	M2tsPath           sql.NullString `json:"m2tsPath"`
	Mp4Path            sql.NullString `json:"mp4Path"`
	Aribb24TxtPath     sql.NullString `json:"aribb24TxtPath"`
	TranscribedTxtPath sql.NullString `json:"transcribedTxtPath"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
}

type Service struct {
	ID                 int64     `json:"id"`
	ServiceID          int32     `json:"serviceId"`
	NetworkID          int32     `json:"networkId"`
	Type               int32     `json:"type"`
	LogoID             int32     `json:"logoID"`
	RemoteControlKeyID int32     `json:"remoteControlKeyId"`
	Name               string    `json:"name"`
	ChannelType        string    `json:"-"`
	Channel            string    `json:"-"`
	HasLogoData        bool      `json:"hasLogoData"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

type TranscribeTask struct {
	ID        int32     `json:"id"`
	TaskID    string    `json:"taskID"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
