
-- +migrate Up
CREATE TABLE IF NOT EXISTS `encode_task` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `task_id` TEXT NOT NULL,
    `status` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `encode_task_task_id_idx` (`task_id`)
);

-- +migrate Down
DROP TABLE `encode_task`;
