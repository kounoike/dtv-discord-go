-- name: GetProgramServiceByServiceID :one
SELECT * FROM `program_service` WHERE `service_id` = ?;

-- name: GetProgramServiceByProgramID :one
SELECT * FROM `program_service` WHERE `program_id` = ?;

-- name: InsertProgramService :exec
INSERT INTO `program_service` (
    `program_id`,
    `service_id`
) VALUES (?, ?);

-- name: GetServiceByProgramID :one
SELECT service.*
FROM `service`
JOIN `program_service` on `program_service`.`service_id` = `service`.`id`
WHERE `program_service`.`program_id` = ?;
