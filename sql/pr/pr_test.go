package pr

import (
	"database/sql"
	"github.com/krlvi/github-devstats/event"
	"github.com/krlvi/github-devstats/sql/schema"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepo_SavePr(t *testing.T) {
	r, err := newRepo()
	assert.NoError(t, err)
	e := fakeEvent()
	err = r.Save(e)
	assert.NoError(t, err)
	ret := r.get(e.Repository, e.PrNumber)
	assert.Equal(t, e, ret)
}

func newRepo() (*Repo, error) {
	db, err := sql.Open("mysql", "devstats:devstats@tcp(127.0.0.1:3306)/devstats?multiStatements=true&parseTime=true")
	if err != nil {
		return nil, err
	}
	r := NewRepo(db)
	_ = schema.MigrateDown(db)
	err = schema.MigrateUp(db)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func fakeEvent() event.Event {
	return event.Event{
		PrNumber:              123,
		Repository:            "foo-bar",
		MergedAt:              time.Unix(1570656320, 0).UTC(),
		TimeToMergeSeconds:    12345,
		BranchAgeSeconds:      22222,
		LinesAdded:            8,
		LinesRemoved:          4,
		FilesChanged:          5,
		CommitsCount:          11,
		CommentsCount:         3,
		AuthorId:              "abc",
		JavaTestFilesModified: 8,
		JavaTestsAdded:        17,
		TimeToApproveSeconds:  123456789,
		ApproverId:            "cba",
		CrossTeam:             true,
		DismissReviewCount:    2,
		ChangesRequestedCount: 1,
	}
}
