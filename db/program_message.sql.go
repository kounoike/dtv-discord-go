// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: program_message.sql

package db

import (
	"context"
)

const getProgramMessageByMessageID = `-- name: GetProgramMessageByMessageID :one
SELECT id, channel_id, message_id, program_id, created_at, updated_at FROM ` + "`" + `program_message` + "`" + ` WHERE ` + "`" + `message_id` + "`" + ` = ?
`

func (q *Queries) GetProgramMessageByMessageID(ctx context.Context, messageID string) (ProgramMessage, error) {
	row := q.db.QueryRowContext(ctx, getProgramMessageByMessageID, messageID)
	var i ProgramMessage
	err := row.Scan(
		&i.ID,
		&i.ChannelID,
		&i.MessageID,
		&i.ProgramID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProgramMessageByProgramID = `-- name: GetProgramMessageByProgramID :one
SELECT id, channel_id, message_id, program_id, created_at, updated_at FROM ` + "`" + `program_message` + "`" + ` WHERE ` + "`" + `program_id` + "`" + ` = ?
`

func (q *Queries) GetProgramMessageByProgramID(ctx context.Context, programID int64) (ProgramMessage, error) {
	row := q.db.QueryRowContext(ctx, getProgramMessageByProgramID, programID)
	var i ProgramMessage
	err := row.Scan(
		&i.ID,
		&i.ChannelID,
		&i.MessageID,
		&i.ProgramID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertProgramMessage = `-- name: InsertProgramMessage :exec
INSERT INTO ` + "`" + `program_message` + "`" + ` (
    ` + "`" + `message_id` + "`" + `,
    ` + "`" + `program_id` + "`" + `,
    ` + "`" + `channel_id` + "`" + `
) VALUES (?, ?, ?)
`

type InsertProgramMessageParams struct {
	MessageID string `json:"messageID"`
	ProgramID int64  `json:"programID"`
	ChannelID string `json:"channelID"`
}

func (q *Queries) InsertProgramMessage(ctx context.Context, arg InsertProgramMessageParams) error {
	_, err := q.db.ExecContext(ctx, insertProgramMessage, arg.MessageID, arg.ProgramID, arg.ChannelID)
	return err
}
