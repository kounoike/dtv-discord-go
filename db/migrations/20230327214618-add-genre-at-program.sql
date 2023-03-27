
-- +migrate Up
ALTER TABLE `program` ADD `genre` TEXT NOT NULL;

-- +migrate Down
ALTER TABLE `program` DROP `genre`;
