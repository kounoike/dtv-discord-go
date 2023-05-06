-- name: InsertAutoSearch :execresult
INSERT INTO `auto_search` (
    `name`,
    `title`,
    `channel`,
    `genre`,
    `kana_search`,
    `fuzzy_search`,
    `regex_search`,
    `record`,
    `thread_id`
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
;

-- name: GetAutoSearch :one
SELECT
    *
FROM
    `auto_search`
WHERE
    `id` = ?
;

-- name: GetAutoSearchByThreadID :one
SELECT
    *
FROM
    `auto_search`
WHERE
    `thread_id` = ?
;

-- name: ListAutoSearch :many
SELECT
    *
FROM
    `auto_search`
;

-- name: DeleteAutoSearch :exec
DELETE FROM `auto_search`
WHERE
    `id` = ?
;
