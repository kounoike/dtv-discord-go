
-- +migrate Up
CREATE TABLE `index_invalid` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `type` TEXT NOT NULL,
    `status` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `index_invalid_type_idx` (`type`)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- +migrate Down
DROP TABLE `index_invalid`;
