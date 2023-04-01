// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: encode_task.sql

package db

import (
	"context"
)

const getEncodeTaskByTaskID = `-- name: GetEncodeTaskByTaskID :one
SELECT id, task_id, status, created_at, updated_at FROM ` + "`" + `encode_task` + "`" + ` WHERE ` + "`" + `task_id` + "`" + ` = ?
`

func (q *Queries) GetEncodeTaskByTaskID(ctx context.Context, taskID string) (EncodeTask, error) {
	row := q.db.QueryRowContext(ctx, getEncodeTaskByTaskID, taskID)
	var i EncodeTask
	err := row.Scan(
		&i.ID,
		&i.TaskID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertEncodeTask = `-- name: InsertEncodeTask :exec
INSERT INTO ` + "`" + `encode_task` + "`" + ` (
    ` + "`" + `task_id` + "`" + `,
    ` + "`" + `status` + "`" + `
) VALUES (?, ?)
`

type InsertEncodeTaskParams struct {
	TaskID string `json:"taskID"`
	Status string `json:"status"`
}

func (q *Queries) InsertEncodeTask(ctx context.Context, arg InsertEncodeTaskParams) error {
	_, err := q.db.ExecContext(ctx, insertEncodeTask, arg.TaskID, arg.Status)
	return err
}
