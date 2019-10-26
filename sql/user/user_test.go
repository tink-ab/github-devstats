package user

import (
	"database/sql"
	"github.com/krlvi/github-devstats/sql/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_Save(t *testing.T) {
	r, err := newRepo()
	assert.NoError(t, err)
	err = r.SaveUser("foo", "Foo Barsson")
	assert.NoError(t, err)
	assert.Equal(t, "Foo Barsson", r.GetName("foo"))
}

func TestRepo_SaveUserTeam(t *testing.T) {
	r, _ := newRepo()
	_ = r.SaveUser("foo", "Foo Barsson")
	err := r.SaveUserTeam("foo", "team-a")
	assert.NoError(t, err)
	err = r.SaveUserTeam("foo", "team-b")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"team-a", "team-b"}, r.getTeamsByUserId("foo"))
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
