CREATE TABLE IF NOT EXISTS pr_commits_by_type
(
    `pr_number` INT REFERENCES prs (`pr_number`),
    `commit_type`      VARCHAR(255),
    `count`     INT,
    PRIMARY KEY (`pr_number`, `commit_type`)
);
