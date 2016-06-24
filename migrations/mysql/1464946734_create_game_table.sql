-- +migrate Up
CREATE TABLE `game` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `team1_id` int NOT NULL,
  `team2_id` int NOT NULL,
  `score1` tinyint NOT NULL,
  `score2` tinyint NOT NULL,
  `modified` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `played_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT `fk_team_1` FOREIGN KEY (`team1_id`) REFERENCES `team` (`id`) ON DELETE RESTRICT,
   CONSTRAINT `fk_team_2` FOREIGN KEY (`team2_id`) REFERENCES `team` (`id`) ON DELETE RESTRICT,
   INDEX `idx_played_at` (`played_at`) 
) CHARACTER SET utf8 COLLATE utf8_general_ci ENGINE=InnoDB;
-- +migrate Down
DROP TABLE `game`;