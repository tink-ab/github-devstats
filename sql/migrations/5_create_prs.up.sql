CREATE TABLE IF NOT EXISTS prs
(
    `pr_number`                INT,
    `repository`               VARCHAR(255),
    `merged_at`                DATETIME NOT NULL,
    `time_to_merge_seconds`    INT,
    `branch_age_seconds`       INT,
    `lines_added`              INT,
    `lines_removed`            INT,
    `files_changed`            INT,
    `commits_count`            INT,
    `comments_count`           INT,
    `author_id`                VARCHAR(255) REFERENCES users (`id`),
    `java_test_files_modified` INT,
    `java_tests_added`         INT,
    `time_to_approve_seconds`  INT,
    `approver_id`              VARCHAR(255) REFERENCES users (`id`),
    `cross_team`               BOOL,
    `dismiss_review_count`     INT,
    `changes_requested_count`  INT,
    PRIMARY KEY (`pr_number`, `repository`),
    INDEX (`merged_at`)
);
