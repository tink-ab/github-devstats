package commit

import (
	"database/sql"
	"github.com/krlvi/github-devstats/event"
	"github.com/krlvi/github-devstats/sql/pr"
	"github.com/krlvi/github-devstats/sql/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepo_SavePrCommitByType(t *testing.T) {
	r, prRepo, err := newRepo()
	assert.NoError(t, err)
	_ = prRepo.Save(event.Event{PrNumber: 333, Repository: "abc"})
	err = r.SavePrCommitByType(333, "abc", "feat", 2)
	err = r.SavePrCommitByType(333, "abc", "fix", 8)
	assert.NoError(t, err)
	assert.Equal(t, 2, r.getCommitTypesByPr(333, "abc")["feat"])
	assert.Equal(t, 8, r.getCommitTypesByPr(333, "abc")["fix"])
}

func TestRepo_SavePrFilesAddedByExt(t *testing.T) {
	r, prRepo, err := newRepo()
	assert.NoError(t, err)
	_ = prRepo.Save(event.Event{PrNumber: 333, Repository: "abc"})
	err = r.SavePrFilesAddedByExt(333, "abc", "py", 3)
	err = r.SavePrFilesAddedByExt(333, "abc", "java", 4)
	assert.NoError(t, err)
	assert.Equal(t, 3, r.getFilesAddedByPr(333, "abc")["py"])
	assert.Equal(t, 4, r.getFilesAddedByPr(333, "abc")["java"])
}

func TestRepo_SavePrFilesModifiedByExt(t *testing.T) {
	r, prRepo, err := newRepo()
	assert.NoError(t, err)
	_ = prRepo.Save(event.Event{PrNumber: 333, Repository: "abc"})
	err = r.SavePrFilesModifiedByExt(333, "abc", "py", 3)
	err = r.SavePrFilesModifiedByExt(333, "abc", "java", 4)
	assert.NoError(t, err)
	assert.Equal(t, 3, r.getFilesModifiedByPr(333, "abc")["py"])
	assert.Equal(t, 4, r.getFilesModifiedByPr(333, "abc")["java"])
}

func newRepo() (*Repo, *pr.Repo, error) {
	db, err := sql.Open("mysql", "devstats:devstats@tcp(127.0.0.1:3306)/devstats?multiStatements=true&parseTime=true")
	if err != nil {
		return nil, nil, err
	}
	_ = schema.MigrateDown(db)
	err = schema.MigrateUp(db)
	if err != nil {
		return nil, nil, err
	}
	return NewRepo(db), pr.NewRepo(db), nil
}
