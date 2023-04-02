-- name: GetComponentVersion :one
SELECT * FROM `component_version` WHERE `component` = ?;

-- name: UpdateComponentVersion :exec
UPDATE `component_version` SET `version` = ? WHERE `component` = ?;

-- name: InsertComponentVersion :exec
INSERT `component_version` (`component`, `version`) VALUES (?, ?);
