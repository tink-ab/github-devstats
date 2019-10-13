package sql

import (
	"database/sql"
	"github.com/krlvi/github-devstats/event"
	"log"
	"sync"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("mysql", "devstats:devstats@tcp(127.0.0.1:3306)/devstats?multiStatements=true&parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ReadAndPersist(repo *Repository, c chan event.Event, wg *sync.WaitGroup) {
	for {
		err := repo.Save(<-c)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}
}
