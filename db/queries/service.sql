-- name: GetService :one
SELECT * FROM `service` WHERE `id` = ?;

-- name: createOrUpdateService :exec
INSERT INTO `service` (
    `id`,
    `service_id`,
    `network_id`,
    `type`,
    `logo_id`,
    `remote_control_key_id`,
    `name`,
    `channel_type`,
    `channel`,
    `has_logo_data`
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    `service_id` = VALUES(`service_id`),
    `network_id` = VALUES(`network_id`),
    `type` = VALUES(`type`),
    `logo_id` = VALUES(`logo_id`),
    `remote_control_key_id` = VALUES(`remote_control_key_id`),
    `channel_type` = VALUES(`channel_type`),
    `channel` = VALUES(`channel`),
    `has_logo_data` = VALUES(`has_logo_data`)
;
