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
