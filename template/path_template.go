package template

import "time"

type PathProgram struct {
	Name      string
	StartTime time.Time
}
type PathService struct {
	Name string
}

type PathTemplateData struct {
	Program  PathProgram `json:"-"`
	Service  PathService `json:"-"`
	Title    string      `json:"title"`
	Subtitle string      `json:"subtitle"`
	Episode  int32       `json:"episode"`
}
