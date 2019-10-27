package user

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

func (r *Repo) SaveUser(id, name string) error {
	_, err := r.db.Exec("INSERT INTO users (`id`, `name`) VALUES (?, ?)", id, name)
	return err
}

func (r *Repo) GetName(id string) (name string) {
	row := r.db.QueryRow("SELECT name FROM users WHERE id = ?", id)
	_ = row.Scan(&name)
	return
}

func (r *Repo) UserExists(id string) (exists bool) {
	row := *r.db.QueryRow("SELECT EXISTS(SELECT id from users WHERE id = ?)", id)
	_ = row.Scan(&exists)
	return exists
}

func (r *Repo) SaveUserTeam(user_id, team_name string) error {
	_, err := r.db.Exec("INSERT INTO user_teams (`user_id`, `team_name`) VALUES (?, ?)", user_id, team_name)
	return err
}

func (r *Repo) RemoveUserTeam(user_id, team_name string) error {
	_, err := r.db.Exec("DELETE FROM user_teams WHERE user_id = ? AND team_name = ?", user_id, team_name)
	return err
}

func (r *Repo) GetTeamsByUserId(user_id string) []string {
	rows, err := r.db.Query("SELECT team_name FROM user_teams WHERE user_id = ?", user_id)
	if err != nil {
		return nil
	}
	var teamNames []string
	for rows.Next() {
		var team_name string
		err = rows.Scan(&team_name)
		if err != nil {
			continue
		}
		teamNames = append(teamNames, team_name)
	}
	return teamNames
}
