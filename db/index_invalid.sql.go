// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: index_invalid.sql

package db

import (
	"context"
)

const setIndexInvalid = `-- name: SetIndexInvalid :exec
INSERT INTO ` + "`" + `index_invalid` + "`" + `
    (` + "`" + `type` + "`" + `, ` + "`" + `status` + "`" + `)
VALUES
    (?, ?)
ON DUPLICATE KEY UPDATE
    ` + "`" + `status` + "`" + ` = VALUES(` + "`" + `status` + "`" + `)
`

type SetIndexInvalidParams struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

func (q *Queries) SetIndexInvalid(ctx context.Context, arg SetIndexInvalidParams) error {
	_, err := q.db.ExecContext(ctx, setIndexInvalid, arg.Type, arg.Status)
	return err
}
