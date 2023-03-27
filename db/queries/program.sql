-- name: GetProgram :one
SELECT * FROM `program` WHERE `id` = ?;

-- name: createProgram :exec
INSERT INTO `program` (
    `id`,
    `json`,
    `event_id`,
    `service_id`,
    `network_id`,
    `start_at`,
    `duration`,
    `is_free`,
    `name`,
    `description`,
    `genre`
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: updateProgram :exec
UPDATE `program` SET
    `json` = ?,
    `event_id` = ?,
    `service_id` = ?,
    `network_id` = ?,
    `start_at` = ?,
    `duration` = ?,
    `is_free` = ?,
    `name` = ?,
    `description` = ?
WHERE id = ?;
