ALTER TABLE pr_events
    ADD COLUMN changes_requested_count INT,
    ADD COLUMN dismiss_review_count INT,
    ADD COLUMN cross_team BOOL,
    ADD COLUMN approver_teams JSON,
    ADD COLUMN approver_name VARCHAR(255),
    ADD COLUMN approver_id VARCHAR(255),
    ADD COLUMN time_to_approve_seconds INT,
    ADD COLUMN java_tests_added INT,
    ADD COLUMN java_test_files_modified INT,
    ADD COLUMN files_modified_by_extension JSON,
    ADD COLUMN files_added_by_extension JSON,
    ADD COLUMN commits_by_type JSON,
    ADD COLUMN author_teams JSON,
    ADD COLUMN author_name VARCHAR(255),
    ADD COLUMN author_id VARCHAR(255),
    ADD COLUMN comments_count INT,
    ADD COLUMN commits_count INT,
    ADD COLUMN files_changed INT,
    ADD COLUMN lines_removed INT,
    ADD COLUMN lines_added INT,
    ADD COLUMN branch_age_seconds INT,
    ADD COLUMN time_to_merge_seconds INT,
    ADD COLUMN merged_at DATETIME NOT NULL;

