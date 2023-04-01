-- name: GetProgramMessageByMessageID :one
SELECT * FROM `program_message` WHERE `message_id` = ?;

-- name: GetProgramMessageByProgramID :one
SELECT * FROM `program_message` WHERE `program_id` = ?;

-- name: InsertProgramMessage :exec
INSERT INTO `program_message` (
    `message_id`,
    `program_id`,
    `channel_id`
) VALUES (?, ?, ?);

