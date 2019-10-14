package commit

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
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

func (r *Repo) SavePrCommitByType(pr_number int, commit_type string, count int) error {
	_, err := r.db.Exec("INSERT INTO pr_commits_by_type (`pr_number`, `commit_type`, `count`) VALUES (?, ?, ?)",
		pr_number, commit_type, count)
	return err
}

func (r *Repo) getCommitTypesByPr(pr_number int) map[string]int {
	rows, err := r.db.Query("SELECT `commit_type`, `count` FROM pr_commits_by_type WHERE pr_number = ?", pr_number)
	if err != nil {
		return nil
	}
	types := map[string]int{}
	for rows.Next() {
		var commit_type string
		var count int
		err = rows.Scan(&commit_type, &count)
		if err != nil {
			continue
		}
		types[commit_type] = count
	}
	return types
}

func (r *Repo) SavePrFilesAddedByExt(pr_number int, ext string, count int) error {
	_, err := r.db.Exec("INSERT INTO pr_files_added_by_ext (`pr_number`, `ext`, `count`) VALUES (?, ?, ?)",
		pr_number, ext, count)
	return err
}

func (r *Repo) getFilesAddedByPr(pr_number int) map[string]int {
	rows, err := r.db.Query("SELECT `ext`, `count` FROM pr_files_added_by_ext WHERE pr_number = ?", pr_number)
	if err != nil {
		return nil
	}
	files := map[string]int{}
	for rows.Next() {
		var ext string
		var count int
		err = rows.Scan(&ext, &count)
		if err != nil {
			continue
		}
		files[ext] = count
	}
	return files
}

func (r *Repo) SavePrFilesModifiedByExt(pr_number int, ext string, count int) error {
	_, err := r.db.Exec("INSERT INTO pr_files_modified_by_ext (`pr_number`, `ext`, `count`) VALUES (?, ?, ?)",
		pr_number, ext, count)
	return err
}

func (r *Repo) getFilesModifiedByPr(pr_number int) map[string]int {
	rows, err := r.db.Query("SELECT `ext`, `count` FROM pr_files_modified_by_ext WHERE pr_number = ?", pr_number)
	if err != nil {
		return nil
	}
	files := map[string]int{}
	for rows.Next() {
		var ext string
		var count int
		err = rows.Scan(&ext, &count)
		if err != nil {
			continue
		}
		files[ext] = count
	}
	return files
}
