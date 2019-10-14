CREATE TABLE IF NOT EXISTS user_teams
(
    `user_id`   VARCHAR(255) REFERENCES users (`id`),
    `team_name` VARCHAR(255),
    PRIMARY KEY (`user_id`, `team_name`)
);

