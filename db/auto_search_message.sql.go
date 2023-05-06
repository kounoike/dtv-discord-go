// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: auto_search_message.sql

package db

import (
	"context"
)

const getAutoSearchMessageByMessageID = `-- name: GetAutoSearchMessageByMessageID :one
SELECT
    id, thread_id, message_id, created_at, updated_at
FROM
    ` + "`" + `auto_search_message` + "`" + `
WHERE
    ` + "`" + `message_id` + "`" + ` = ?
`

func (q *Queries) GetAutoSearchMessageByMessageID(ctx context.Context, messageID string) (AutoSearchMessage, error) {
	row := q.db.QueryRowContext(ctx, getAutoSearchMessageByMessageID, messageID)
	var i AutoSearchMessage
	err := row.Scan(
		&i.ID,
		&i.ThreadID,
		&i.MessageID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertAutoSearchMessage = `-- name: InsertAutoSearchMessage :exec
INSERT INTO ` + "`" + `auto_search_message` + "`" + ` (
    ` + "`" + `thread_id` + "`" + `,
    ` + "`" + `message_id` + "`" + `
) VALUES (?, ?)
`

type InsertAutoSearchMessageParams struct {
	ThreadID  string `json:"threadID"`
	MessageID string `json:"messageID"`
}

func (q *Queries) InsertAutoSearchMessage(ctx context.Context, arg InsertAutoSearchMessageParams) error {
	_, err := q.db.ExecContext(ctx, insertAutoSearchMessage, arg.ThreadID, arg.MessageID)
	return err
}
