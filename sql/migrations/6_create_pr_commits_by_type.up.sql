CREATE TABLE IF NOT EXISTS pr_commits_by_type
(
    `pr_number` INT REFERENCES prs (`pr_number`),
    `repo` VARCHAR(255),
    `commit_type`      VARCHAR(255),
    `count`     INT,
    PRIMARY KEY (`pr_number`, `repo`, `commit_type`)
);
