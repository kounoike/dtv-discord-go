-- name: SetIndexInvalid :exec
INSERT INTO `index_invalid`
    (`type`, `status`)
VALUES
    (?, ?)
ON DUPLICATE KEY UPDATE
    `status` = VALUES(`status`);
