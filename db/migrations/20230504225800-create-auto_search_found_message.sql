CREATE TABLE IF NOT EXISTS `auto_search_found_message` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `thread_id` TEXT NOT NULL,
    `program_id` BIGINT UNSIGNED NOT NULL,
    `message_id` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `auto_search_found_message_message_id_idx` (`message_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
