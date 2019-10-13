package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/krlvi/github-devstats/event"
	"log"
	"os"
	"reflect"
	"strings"
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
	approverTeams, err := json.Marshal(e.ApproverTeams)
	if err != nil {
		return err
	}
	log.Println("event at", e.MergedAt, "saving repo", e.Repository, "pr", e.PrNumber)
	_, err = r.db.Exec("INSERT INTO pr_events ("+tableColumns()+") VALUES ("+tableColumnPlaceholders()+")",
		e.PrNumber,
		e.Repository,
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
		approverTeams,
		e.CrossTeam,
		e.DismissReviewCount,
		e.ChangesRequestedCount,
	)
	if err != nil {
		return err
	}
	return nil
}

// Test only
func (r *Repository) get(repository string, pr_number int) event.Event {
	row := r.db.QueryRow("SELECT "+tableColumns()+" FROM pr_events WHERE repository = ? AND pr_number = ?",
		repository, pr_number)
	e := event.Event{}
	var authorTeams []byte
	var commitsByType []byte
	var filesAddedByExtension []byte
	var filesModifiedByExtension []byte
	var approverTeams []byte
	_ = row.Scan(
		&e.PrNumber,
		&e.Repository,
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
		&approverTeams,
		&e.CrossTeam,
		&e.DismissReviewCount,
		&e.ChangesRequestedCount,
	)
	_ = json.Unmarshal(authorTeams, &e.AuthorTeams)
	_ = json.Unmarshal(commitsByType, &e.CommitsByType)
	_ = json.Unmarshal(filesAddedByExtension, &e.FilesAddedByExtension)
	_ = json.Unmarshal(filesModifiedByExtension, &e.FilesModifiedByExtension)
	_ = json.Unmarshal(approverTeams, &e.ApproverTeams)
	return e
}

func tableColumns() string {
	val := reflect.ValueOf(event.Event{})
	var sb strings.Builder
	for i := 0; i < val.Type().NumField(); i++ {
		sb.WriteRune('`')
		sb.WriteString(val.Type().Field(i).Tag.Get("json"))
		sb.WriteRune('`')
		if i < val.Type().NumField()-1 {
			sb.WriteRune(',')
		}
	}
	return sb.String()
}

func tableColumnPlaceholders() string {
	val := reflect.ValueOf(event.Event{})
	var sb strings.Builder
	for i := 0; i < val.Type().NumField(); i++ {
		sb.WriteRune('?')
		if i < val.Type().NumField()-1 {
			sb.WriteRune(',')
		}
	}
	return sb.String()
}
