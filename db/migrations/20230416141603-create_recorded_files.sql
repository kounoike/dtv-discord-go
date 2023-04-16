
-- +migrate Up
CREATE TABLE IF NOT EXISTS `recorded_files` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `program_id` BIGINT UNSIGNED NOT NULL,
  `m2ts_path` text NULL DEFAULT NULL,
  `mp4_path` text NULL DEFAULT NULL,
  `aribb24_txt_path` text NULL DEFAULT NULL,
  `transcribed_txt_path` text NULL DEFAULT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- +migrate Down
DROP TABLE `recorded_files`;
