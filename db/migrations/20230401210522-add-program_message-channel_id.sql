
-- +migrate Up
ALTER TABLE `program_message` ADD `channel_id` TEXT NOT NULL;

-- +migrate Down
ALTER TABLE `program_message` DROP `channel_id`;
