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
SELECT service.id, service.service_id, service.network_id, service.type, service.logo_id, service.remote_control_key_id, service.name, service.channel_type, service.channel, service.has_logo_data, service.created_at, service.updated_at
FROM `service`
JOIN `program_service` on `program_service`.`service_id` = `service`.`id`
WHERE `program_service`.`program_id` = ?;
