package sql

import (
	"database/sql"
	"github.com/krlvi/github-devstats/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepository_Save(t *testing.T) {
	r := NewRepo(t)
	defer r.migrateDown()
	e := FakeEvent()
	assert.NoError(t, r.Save(e))
	persisted := r.get(e.Repository, e.PrNumber)
	assert.Equal(t, e, persisted)
}

func TestRepository_Migrations(t *testing.T) {
	r := NewRepo(t)
	err := r.migrateDown()
	if err != nil {
		t.Log(err)
		panic(err)
	}
}

func NewRepo(t *testing.T) *Repository {
	db, err := sql.Open("mysql", "devstats:devstats@tcp(127.0.0.1:3306)/devstats?multiStatements=true")
	if err != nil {
		t.Log(err)
		panic(err)
	}
	r, err := NewRepository(db)
	if err != nil {
		t.Log(err)
		panic(err)
	}
	err = r.MigrateUp()
	if err != nil {
		t.Log(err)
		panic(err)
	}
	return r
}

func FakeEvent() event.Event {
	return event.Event{
		PrNumber:           123,
		Repository:         "foo-bar",
		MergedAt:           time.Unix(1570656320, 0),
		TimeToMergeSeconds: 12345,
		BranchAgeSeconds:   22222,
		LinesAdded:         8,
		LinesRemoved:       4,
		FilesChanged:       5,
		CommitsCount:       11,
		CommentsCount:      3,
		AuthorId:           "abc",
		AuthorName:         "Foo Barsson",
	}
}
