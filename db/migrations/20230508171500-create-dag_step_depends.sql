CREATE TABLE IF NOT EXISTS `dag_step_depends` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `dag_id` INT UNSIGNED NOT NULL,
    `dag_step_id` INT UNSIGNED NOT NULL,
    `dag_step_depend_id` INT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT FOREIGN KEY (`dag_id`) REFERENCES `dag`(`id`),
    CONSTRAINT FOREIGN KEY (`dag_step_id`) REFERENCES `dag_step`(`id`),
    CONSTRAINT FOREIGN KEY (`dag_step_depend_id`) REFERENCES `dag_step`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
