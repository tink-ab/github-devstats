package sql

import (
	"database/sql"
	"github.com/krlvi/github-devstats/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepository_Save(t *testing.T) {
	r, err := NewRepo()
	assert.NoError(t, err)
	e := FakeEvent()
	assert.NoError(t, r.Save(e))
	persisted := r.get(e.Repository, e.PrNumber)
	assert.Equal(t, e, persisted)
}

func TestRepository_Migrations(t *testing.T) {
	_, err := NewRepo()
	assert.NoError(t, err)
}

func NewRepo() (*Repository, error) {
	db, err := sql.Open("mysql", "devstats:devstats@tcp(127.0.0.1:3306)/devstats?multiStatements=true&parseTime=true")
	if err != nil {
		return nil, err
	}
	r, err := NewRepository(db)
	if err != nil {
		return nil, err
	}
	_ = r.migrateDown()
	err = r.MigrateUp()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func FakeEvent() event.Event {
	return event.Event{
		PrNumber:                 123,
		Repository:               "foo-bar",
		MergedAt:                 time.Unix(1570656320, 0).UTC(),
		TimeToMergeSeconds:       12345,
		BranchAgeSeconds:         22222,
		LinesAdded:               8,
		LinesRemoved:             4,
		FilesChanged:             5,
		CommitsCount:             11,
		CommentsCount:            3,
		AuthorId:                 "abc",
		AuthorName:               "Foo Barsson",
		AuthorTeams:              []string{"foo-team", "bar-squad"},
		CommitsByType:            map[string]int{"fix": 1, "feat": 2},
		FilesAddedByExtension:    map[string]int{"java": 4, "BUILD": 1},
		FilesModifiedByExtension: map[string]int{"go": 4, "py": 1},
		JavaTestFilesModified:    8,
		JavaTestsAdded:           17,
		TimeToApproveSeconds:     123456789,
		ApproverId:               "cba",
		ApproverName:             "Bar Foosson",
	}
}
