package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/krlvi/github-devstats/event"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	db       *sql.DB
	migrator *migrate.Migrate
}

func NewRepository(db *sql.DB) (*Repository, error) {
	migrator, err := newMigrator(db)
	if err != nil {
		return nil, err
	}
	repo := &Repository{
		db:       db,
		migrator: migrator,
	}
	return repo, nil
}

func (r *Repository) MigrateUp() error {
	return r.migrator.Up()
}

// Test only
func (r *Repository) migrateDown() error {
	return r.migrator.Down()
}

func newMigrator(db *sql.DB) (*migrate.Migrate, error) {
	var migrationsDir string
	if srcdir, ok := os.LookupEnv("TEST_SRCDIR"); ok {
		migrationsDir = srcdir + "/__main__/sql/migrations"
	} else if wd, ok := os.LookupEnv("BUILD_WORKING_DIRECTORY"); ok {
		migrationsDir = wd + "/sql/migrations"
	}
	log.Printf("Loading migrations from: %s", migrationsDir)
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	return migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"mysql",
		driver,
	)
}

func (r *Repository) Save(e event.Event) error {
	authorTeams, err := json.Marshal(e.AuthorTeams)
	commitsByType, err := json.Marshal(e.CommitsByType)
	filesAddedByExtension, err := json.Marshal(e.FilesAddedByExtension)
	filesModifiedByExtension, err := json.Marshal(e.FilesModifiedByExtension)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("INSERT INTO pr_events ("+
		"`repository`,"+
		"`pr_number`,"+
		"`merged_at`,"+
		"`time_to_merge_seconds`,"+
		"`branch_age_seconds`,"+
		"`lines_added`,"+
		"`lines_removed`,"+
		"`files_changed`,"+
		"`commits_count`,"+
		"`comments_count`,"+
		"`author_id`,"+
		"`author_name`,"+
		"`author_teams`,"+
		"`commits_by_type`,"+
		"`files_added_by_extension`,"+
		"`files_modified_by_extension`,"+
		"`java_test_files_modified`,"+
		"`java_tests_added`,"+
		"`time_to_approve_seconds`,"+
		"`approver_id`,"+
		"`approver_name`"+
		") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		e.Repository,
		e.PrNumber,
		e.MergedAt,
		e.TimeToMergeSeconds,
		e.BranchAgeSeconds,
		e.LinesAdded,
		e.LinesRemoved,
		e.FilesChanged,
		e.CommitsCount,
		e.CommentsCount,
		e.AuthorId,
		e.AuthorName,
		authorTeams,
		commitsByType,
		filesAddedByExtension,
		filesModifiedByExtension,
		e.JavaTestFilesModified,
		e.JavaTestsAdded,
		e.TimeToApproveSeconds,
		e.ApproverId,
		e.ApproverName,
		)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) get(repository string, pr_number int) event.Event {
	row := r.db.QueryRow("SELECT"+
		"`repository`,"+
		"`pr_number`,"+
		"`merged_at`,"+
		"`time_to_merge_seconds`,"+
		"`branch_age_seconds`,"+
		"`lines_added`,"+
		"`lines_removed`,"+
		"`files_changed`,"+
		"`commits_count`,"+
		"`comments_count`,"+
		"`author_id`,"+
		"`author_name`,"+
		"`author_teams`,"+
		"`commits_by_type`,"+
		"`files_added_by_extension`,"+
		"`files_modified_by_extension`,"+
		"`java_test_files_modified`,"+
		"`java_tests_added`,"+
		"`time_to_approve_seconds`,"+
		"`approver_id`,"+
		"`approver_name`"+
		" FROM pr_events WHERE repository = ? AND pr_number = ?", repository, pr_number)
	e := event.Event{}
	var authorTeams []byte
	var commitsByType []byte
	var filesAddedByExtension []byte
	var filesModifiedByExtension []byte
	_ = row.Scan(
		&e.Repository,
		&e.PrNumber,
		&e.MergedAt,
		&e.TimeToMergeSeconds,
		&e.BranchAgeSeconds,
		&e.LinesAdded,
		&e.LinesRemoved,
		&e.FilesChanged,
		&e.CommitsCount,
		&e.CommentsCount,
		&e.AuthorId,
		&e.AuthorName,
		&authorTeams,
		&commitsByType,
		&filesAddedByExtension,
		&filesModifiedByExtension,
		&e.JavaTestFilesModified,
		&e.JavaTestsAdded,
		&e.TimeToApproveSeconds,
		&e.ApproverId,
		&e.ApproverName,
	)
	_ = json.Unmarshal(authorTeams, &e.AuthorTeams)
	_ = json.Unmarshal(commitsByType, &e.CommitsByType)
	_ = json.Unmarshal(filesAddedByExtension, &e.FilesAddedByExtension)
	_ = json.Unmarshal(filesModifiedByExtension, &e.FilesModifiedByExtension)
	return e
}
