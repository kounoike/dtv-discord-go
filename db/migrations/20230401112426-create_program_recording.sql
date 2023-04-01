
-- +migrate Up
CREATE TABLE IF NOT EXISTS `program_recording` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `program_id` BIGINT UNSIGNED NOT NULL,
    `content_path` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `program_recording_program_id_idx` (`program_id`)
);

-- +migrate Down
DROP TABLE `program_recording`;
