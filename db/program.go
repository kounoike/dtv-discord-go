package db

import (
	"context"
	"encoding/json"
	"time"
)

func (p *Program) UnmarshalJSON(b []byte) error {
	type program Program
	var pp program
	err := json.Unmarshal(b, &pp)
	if err != nil {
		return err
	}
	*p = (Program)(pp)
	p.Json = b
	return nil
}

func (q *Queries) InsertProgram(ctx context.Context, p Program) error {
	args := createProgramParams{
		ID:          p.ID,
		Json:        p.Json,
		EventID:     p.EventID,
		ServiceID:   p.ServiceID,
		NetworkID:   p.NetworkID,
		StartAt:     p.StartAt,
		Duration:    p.Duration,
		IsFree:      p.IsFree,
		Name:        p.Name,
		Description: p.Description,
	}
	return q.createProgram(ctx, args)
}

func (q *Queries) UpdateProgram(ctx context.Context, p Program) error {
	args := updateProgramParams{
		ID:          p.ID,
		Json:        p.Json,
		EventID:     p.EventID,
		ServiceID:   p.ServiceID,
		NetworkID:   p.NetworkID,
		StartAt:     p.StartAt,
		Duration:    p.Duration,
		IsFree:      p.IsFree,
		Name:        p.Name,
		Description: p.Description,
	}
	return q.updateProgram(ctx, args)
}

func (p *Program) StartTime() time.Time {
	return time.Unix(p.StartAt/1000, (p.StartAt%1000)*1000)
}
