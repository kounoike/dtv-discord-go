-- name: InsertRecordedFiles :exec
INSERT INTO `recorded_files` (
    `program_id`,
    `m2ts_path`
) VALUES (?, ?);

-- name: UpdateRecordedFilesMp4 :exec
UPDATE `recorded_files` SET
    `mp4_path` = ?
WHERE `program_id` = ?;

-- name: UpdateRecordedFilesAribb24Txt :exec
UPDATE `recorded_files` SET
    `aribb24_txt_path` = ?
WHERE `program_id` = ?;

-- name: UpdateRecordedFilesTranscribedTxt :exec
UPDATE `recorded_files` SET
    `transcribed_txt_path` = ?
WHERE `program_id` = ?;

-- name: ListRecordedFiles :many
SELECT
    `recorded_files`.`program_id`,
    `recorded_files`.`m2ts_path`,
    `recorded_files`.`mp4_path`,
    `recorded_files`.`aribb24_txt_path`,
    `recorded_files`.`transcribed_txt_path`,
    `program`.`json`,
    `program`.`start_at`,
    `program`.`duration`,
    `program`.`name`,
    `program`.`description`,
    `program`.`genre`,
    `service`.name AS service_name
FROM `recorded_files`
JOIN `program` ON `program`.`id` = `recorded_files`.`program_id`
JOIN `service` ON `program`.`service_id` = `service`.`service_id` AND `program`.`network_id` = `service`.`network_id`
;

