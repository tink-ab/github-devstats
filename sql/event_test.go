package sql

import (
	"database/sql"
	"github.com/krlvi/github-devstats/event"
	"testing"
)

func TestRepository_Save(t *testing.T) {
	r := NewRepo(t)
	defer r.migrateDown()
	err := r.Save(FakeEvent())
	if err != nil {
		t.Log(err)
		t.Fail()
	}
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
		PrNumber:   123,
		Repository: "foo-bar",
	}
}
