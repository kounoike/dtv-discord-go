-- name: GetProgram :one
SELECT * FROM `program` WHERE `id` = ?;

-- name: createProgram :exec
INSERT INTO `program` (
    `id`,
    `json`,
    `event_id`,
    `service_id`,
    `network_id`,
    `start_at`,
    `duration`,
    `is_free`,
    `name`,
    `description`,
    `genre`
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: updateProgram :exec
UPDATE `program` SET
    `json` = ?,
    `event_id` = ?,
    `service_id` = ?,
    `network_id` = ?,
    `start_at` = ?,
    `duration` = ?,
    `is_free` = ?,
    `name` = ?,
    `description` = ?
WHERE id = ?;

-- name: ListProgramWithMessageAndServiceName :many
SELECT
    `program`.`id` AS `program_id`,
    `program`.`json`,
    `program`.`event_id`,
    `program`.`service_id`,
    `program`.`network_id`,
    `program`.`start_at`,
    `program`.`duration`,
    `program`.`is_free`,
    `program`.`name`,
    `program`.`description`,
    `program`.`genre`,
    `program_message`.`channel_id`,
    `program_message`.`message_id`,
    `service`.`name` AS `service_name`
FROM `program`
JOIN `program_message` ON `program_message`.`program_id` = `program`.`id`
JOIN `service` ON `program`.`service_id` = `service`.`service_id` AND `program`.`network_id` = `service`.`network_id`
;
