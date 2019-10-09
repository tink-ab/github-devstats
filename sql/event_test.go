package sql

import (
	"database/sql"
	"testing"
)

func TestRepository_Migrations(t *testing.T) {
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
	_ = r.migrateDown()
	err = r.MigrateUp()
	if err != nil {
		t.Log(err)
		panic(err)
	}
	err = r.migrateDown()
	if err != nil {
		t.Log(err)
		panic(err)
	}
}
