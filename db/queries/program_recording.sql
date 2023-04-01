-- name: GetProgramRecording :one
SELECT * FROM program_recording WHERE id = ?;

-- name: GetProgramRecordingByProgramId :one
SELECT * FROM program_recording WHERE program_id = ?;

-- name: InsertProgramRecording :exec
INSERT INTO program_recording(
    program_id,
    content_path
) VALUES (
    ?,
    ?
);

-- name: DeleteProgramRecordingByProgramId :exec
DELETE FROM program_recording WHERE program_id = ?;
