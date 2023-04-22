CREATE TABLE IF NOT EXISTS `program_message` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `channel_id` TEXT NOT NULL,
    `message_id` TEXT NOT NULL,
    `program_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `program_message_message_id_idx` (`message_id`),
    INDEX `program_message_program_id_idx` (`program_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
