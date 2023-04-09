-- name: GetTranscribeTaskByTaskID :one
SELECT * FROM `transcribe_task` WHERE `task_id` = ?;

-- name: InsertTranscribeTask :exec
INSERT INTO `transcribe_task` (
    `task_id`,
    `status`
) VALUES (?, ?);
