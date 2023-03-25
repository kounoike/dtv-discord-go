
-- +migrate Up
CREATE TABLE IF NOT EXISTS `service` (
    `id` BIGINT UNSIGNED NOT NULL,
    `service_id` INT UNSIGNED NOT NULL,
    `network_id` INT UNSIGNED NOT NULL,
    `type` INT NOT NULL,
    `logo_id` INT NOT NULL,
    `remote_control_key_id` INT NOT NULL,
    `name` TEXT NOT NULL,
    `channel_type`  TEXT NOT NULL,
    `channel` TEXT NOT NULL,
    `has_logo_data` BOOLEAN NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

-- +migrate Down
DROP TABLE `service`;
