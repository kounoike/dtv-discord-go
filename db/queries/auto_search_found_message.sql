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

-- name: CountAutoSearchFoundMessagesWithRecordByProgramID :one
SELECT 
    count(*)
FROM
    `auto_search_found_message`
JOIN
    `auto_search` ON `auto_search_found_message`.`thread_id` = `auto_search`.`thread_id`
WHERE
    `auto_search`.`record` = 1
    AND `program_id` = ?
;

-- name: DeleteAutoSearchFoundMessagesByThreadID :exec
DELETE FROM `auto_search_found_message`
WHERE
    `thread_id` = ?
;
