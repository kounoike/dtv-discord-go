-- name: InsertAutoSearchFoundMessage :exec
INSERT INTO `auto_search_found_message` (
    `thread_id`,
    `program_id`,
    `message_id`
) VALUES (?, ?, ?);

-- name: ListAutoSearchFoundMessages :many
SELECT 
    *
FROM
    `auto_search_found_message`
WHERE
    `thread_id` = ?;

-- name: CountAutoSearchFoundMessagesByProgramID :one
SELECT 
    count(*)
FROM
    `auto_search_found_message`
WHERE
    `program_id` = ?
;

-- name: DeleteAutoSearchFoundMessagesByThreadID :exec
DELETE FROM `auto_search_found_message`
WHERE
    `thread_id` = ?
;
