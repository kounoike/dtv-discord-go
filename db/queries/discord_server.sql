-- name: InsertServerId :exec
INSERT INTO discord_server (server_id)
VALUES (?)
ON DUPLICATE KEY UPDATE
    server_id = VALUES(server_id);

