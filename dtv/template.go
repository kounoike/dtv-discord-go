package dtv

import "time"

type PathProgram struct {
	Name      string
	StartTime time.Time
}
type PathService struct {
	Name string
}

type PathTemplateData struct {
	Program PathProgram
	Service PathService
}
