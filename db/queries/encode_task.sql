-- name: GetEncodeTaskByTaskID :one
SELECT * FROM `encode_task` WHERE `task_id` = ?;

-- name: InsertEncodeTask :exec
INSERT INTO `encode_task` (
    `task_id`,
    `status`
) VALUES (?, ?);
