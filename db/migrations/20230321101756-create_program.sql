
-- +migrate Up
CREATE TABLE IF NOT EXISTS `program` (
    `id` BIGINT UNSIGNED NOT NULL,
    `json` JSON NOT NULL,
    `event_id` INT UNSIGNED NOT NULL,
    `service_id` INT UNSIGNED NOT NULL,
    `network_id` INT UNSIGNED NOT NULL,
    `start_at` BIGINT UNSIGNED NOT NULL,
    `duration` INT UNSIGNED NOT NULL,
    `is_free` BOOLEAN NOT NULL,
    `name` TEXT NOT NULL,
    `description` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

-- +migrate Down
DROP TABLE `program`;
