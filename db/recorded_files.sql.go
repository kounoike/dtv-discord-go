// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: recorded_files.sql

package db

import (
	"context"
	"database/sql"
)

const insertRecordedFiles = `-- name: InsertRecordedFiles :exec
INSERT INTO ` + "`" + `recorded_files` + "`" + ` (
    ` + "`" + `program_id` + "`" + `,
    ` + "`" + `m2ts_path` + "`" + `
) VALUES (?, ?)
`

type InsertRecordedFilesParams struct {
	ProgramID int64          `json:"programID"`
	M2tsPath  sql.NullString `json:"m2tsPath"`
}

func (q *Queries) InsertRecordedFiles(ctx context.Context, arg InsertRecordedFilesParams) error {
	_, err := q.db.ExecContext(ctx, insertRecordedFiles, arg.ProgramID, arg.M2tsPath)
	return err
}

const updateRecordedFilesAribb24Txt = `-- name: UpdateRecordedFilesAribb24Txt :exec
UPDATE ` + "`" + `recorded_files` + "`" + ` SET
    ` + "`" + `aribb24_txt_path` + "`" + ` = ?
WHERE ` + "`" + `program_id` + "`" + ` = ?
`

type UpdateRecordedFilesAribb24TxtParams struct {
	Aribb24TxtPath sql.NullString `json:"aribb24TxtPath"`
	ProgramID      int64          `json:"programID"`
}

func (q *Queries) UpdateRecordedFilesAribb24Txt(ctx context.Context, arg UpdateRecordedFilesAribb24TxtParams) error {
	_, err := q.db.ExecContext(ctx, updateRecordedFilesAribb24Txt, arg.Aribb24TxtPath, arg.ProgramID)
	return err
}

const updateRecordedFilesMp4 = `-- name: UpdateRecordedFilesMp4 :exec
UPDATE ` + "`" + `recorded_files` + "`" + ` SET
    ` + "`" + `mp4_path` + "`" + ` = ?
WHERE ` + "`" + `program_id` + "`" + ` = ?
`

type UpdateRecordedFilesMp4Params struct {
	Mp4Path   sql.NullString `json:"mp4Path"`
	ProgramID int64          `json:"programID"`
}

func (q *Queries) UpdateRecordedFilesMp4(ctx context.Context, arg UpdateRecordedFilesMp4Params) error {
	_, err := q.db.ExecContext(ctx, updateRecordedFilesMp4, arg.Mp4Path, arg.ProgramID)
	return err
}

const updateRecordedFilesTranscribedTxt = `-- name: UpdateRecordedFilesTranscribedTxt :exec
UPDATE ` + "`" + `recorded_files` + "`" + ` SET
    ` + "`" + `transcribed_txt_path` + "`" + ` = ?
WHERE ` + "`" + `program_id` + "`" + ` = ?
`

type UpdateRecordedFilesTranscribedTxtParams struct {
	TranscribedTxtPath sql.NullString `json:"transcribedTxtPath"`
	ProgramID          int64          `json:"programID"`
}

func (q *Queries) UpdateRecordedFilesTranscribedTxt(ctx context.Context, arg UpdateRecordedFilesTranscribedTxtParams) error {
	_, err := q.db.ExecContext(ctx, updateRecordedFilesTranscribedTxt, arg.TranscribedTxtPath, arg.ProgramID)
	return err
}
