package db

import (
	"context"
	"encoding/json"
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

func (p *Program) InsertDb(ctx context.Context, q Queries) error {
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

func (p *Program) UpdateDb(ctx context.Context, q Queries) error {
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
