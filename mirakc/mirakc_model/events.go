package mirakc_model

type ProgramsUpdatedEventData struct {
	ServiceId uint `json:"serviceId"`
}

type ProgramEventData struct {
	ProgramId int64 `json:"programId"`
}
