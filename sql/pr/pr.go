package pr

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/krlvi/github-devstats/event"
)

type Repo struct {
	db       *sql.DB
	migrator *migrate.Migrate
}

func NewRepo(db *sql.DB) *Repo {
	repo := &Repo{
		db: db,
	}
	return repo
}

func (r *Repo) PrExists(repo string, prNum int) (exists bool) {
	row := *r.db.QueryRow("SELECT EXISTS(SELECT repository, pr_number from prs WHERE repository = ? AND pr_number = ?)", repo, prNum)
	_ = row.Scan(&exists)
	return exists
}

func (r *Repo) Save(e event.Event) error {
	_, err := r.db.Exec("INSERT INTO prs ("+
		"`pr_number`,"+
		"`repository`,"+
		"`merged_at`,"+
		"`time_to_merge_seconds`,"+
		"`branch_age_seconds`,"+
		"`lines_added`,"+
		"`lines_removed`,"+
		"`files_changed`,"+
		"`commits_count`,"+
		"`comments_count`,"+
		"`author_id`,"+
		"`java_test_files_modified`,"+
		"`java_tests_added`,"+
		"`go_test_files_modified`,"+
		"`go_tests_added`,"+
		"`time_to_approve_seconds`,"+
		"`approver_id`,"+
		"`cross_team`,"+
		"`dismiss_review_count`,"+
		"`changes_requested_count`"+
		") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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
		e.JavaTestFilesModified,
		e.JavaTestsAdded,
		e.GoTestFilesModified,
		e.GoTestsAdded,
		e.TimeToApproveSeconds,
		e.ApproverId,
		e.CrossTeam,
		e.DismissReviewCount,
		e.ChangesRequestedCount,
	)
	return err
}

func (r *Repo) get(repository string, pr_number int) event.Event {
	row := r.db.QueryRow("SELECT "+
		"`pr_number`,"+
		"`repository`,"+
		"`merged_at`,"+
		"`time_to_merge_seconds`,"+
		"`branch_age_seconds`,"+
		"`lines_added`,"+
		"`lines_removed`,"+
		"`files_changed`,"+
		"`commits_count`,"+
		"`comments_count`,"+
		"`author_id`,"+
		"`java_test_files_modified`,"+
		"`java_tests_added`,"+
		"`go_test_files_modified`,"+
		"`go_tests_added`,"+
		"`time_to_approve_seconds`,"+
		"`approver_id`,"+
		"`cross_team`,"+
		"`dismiss_review_count`,"+
		"`changes_requested_count`"+
		" FROM prs WHERE repository = ? AND pr_number = ?",
		repository, pr_number)
	e := event.Event{}
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
		&e.JavaTestFilesModified,
		&e.JavaTestsAdded,
		&e.GoTestFilesModified,
		&e.GoTestsAdded,
		&e.TimeToApproveSeconds,
		&e.ApproverId,
		&e.CrossTeam,
		&e.DismissReviewCount,
		&e.ChangesRequestedCount,
	)
	return e
}
