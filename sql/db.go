package sql

import (
	"database/sql"
	"github.com/krlvi/github-devstats/event"
	"log"
)

func Read(c <-chan event.Event) error {
	db, err := sql.Open("mysql", "devstats:devstats@tcp(127.0.0.1:3306)/devstats?multiStatements=true&parseTime=true")
	if err != nil {
		return err
	}
	defer db.Close()
	repo, err := NewRepository(db)
	if err != nil {
		return err
	}
	err = repo.Save(<-c)
	if err != nil {
		log.Println(err)
	}
	return nil
}
