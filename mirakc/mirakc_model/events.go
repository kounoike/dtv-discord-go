package mirakc_model

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
)

type ProgramsUpdatedEventData struct {
	ServiceId uint `json:"serviceId"`
}

type ProgramEventData struct {
	ProgramId int64 `json:"programId"`
}

type Reason struct {
	Type     string      `json:"type"`
	Message  null.String `json:"message"`
	OsError  null.Int    `json:"osError"`
	ExitCode null.Int    `json:"exitCode"`
}

type RecordingFailedEventData struct {
	ProgramId int64  `json:"programId"`
	Reason    Reason `json:"reason"`
}

func (d *RecordingFailedEventData) ToString() string {
	ret := fmt.Sprintf("Failed:%d reason:%s", d.ProgramId, d.Reason.Type)
	if d.Reason.Message.Valid {
		ret += " " + d.Reason.Message.String
	}
	if d.Reason.OsError.Valid {
		ret += fmt.Sprintf(" OsError:%d", d.Reason.OsError.Int64)
	}
	if d.Reason.ExitCode.Valid {
		ret += fmt.Sprintf(" ExitCode:%d", d.Reason.ExitCode.Int64)
	}
	return ret
}
