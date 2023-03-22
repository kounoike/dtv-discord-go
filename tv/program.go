package tv

import "time"

type Genre struct {
	Lv1 uint `json:"lv1"`
	Lv2 uint `json:"lv2"`
	Un1 uint `json:"un1"`
	Un2 uint `json:"un2"`
}

type Extend struct {
	Key   string
	Value string
}

type Program struct {
	ID          uint64   `json:"id"`
	EventID     uint     `json:"eventId"`
	ServiceID   uint     `json:"serviceId"`
	NetworkID   uint     `json:"networkId"`
	StartAt     uint64   `json:"startAt"`
	Duration    uint     `json:"duration"`
	IsFree      bool     `json:"isFree"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Genres      []Genre  `json:"genres"`
	Extends     []Extend `json:"-"`
}

func (p *Program) GetStartTime() time.Time {
	return time.Unix(int64(p.StartAt)/1000, (int64(p.StartAt)%1000)*1000)
}

func (p *Program) GetEndTime() time.Time {
	endAt := p.StartAt + uint64(p.Duration)
	return time.Unix(int64(endAt)/1000, (int64(endAt)%1000)*1000)
}
