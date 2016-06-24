-- +migrate Up
ALTER TABLE `team`
ADD `player1_id` int NULL,
ADD `player2_id` int NULL AFTER `player1_id`;
ALTER TABLE `team`
ADD CONSTRAINT `fk_player_1` FOREIGN KEY (`player1_id`) REFERENCES `player` (`id`) ON DELETE SET NULL,
ADD CONSTRAINT `fk_player_2` FOREIGN KEY (`player2_id`) REFERENCES `player` (`id`) ON DELETE SET NULL;

-- +migrate Down
ALTER TABLE `team`
DROP FOREIGN KEY `fk_player_1`,
DROP FOREIGN KEY `fk_player_2`;
ALTER TABLE `team`
DROP `player1_id`,
DROP `player2_id`;


