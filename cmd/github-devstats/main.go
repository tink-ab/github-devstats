package main

import (
	"github.com/google/go-github/github"
	access "github.com/krlvi/github-devstats/sql"
	"github.com/krlvi/github-devstats/sql/user"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/krlvi/github-devstats/client"
	"github.com/krlvi/github-devstats/event"
)

func main() {
	if len(os.Args) < 3 {
		printUsageAndExit()
	}
	org := os.Args[1]
	accessToken := os.Args[2]
	if len(org) <= 0 || len(accessToken) <= 0 {
		printUsageAndExit()
	}
	daysAgo := 0
	if len(os.Args) > 3 {
		daysAgo, _ = strconv.Atoi(os.Args[3])
	}

	refresh := false
	if len(os.Args) > 4 && os.Args[4] == "refresh" {
		refresh = true
	}

	log.Println("fetching github teams and their members for", org)
	c := client.NewClient(org, accessToken)

	date := time.Now().AddDate(0, 0, daysAgo*-1)
	log.Println("fetching merged pull requests for", org, "on date", date.Format("2006-01-02"))
	prIssues, err := c.GetAllMergedPRIssues(date)
	if err != nil {
		log.Panicln("could not fetch pull requests:", err)
	}

	err = processIntoDB(c, prIssues, refresh)
	if err != nil {
		log.Println(err)
	}
}

func processIntoDB(c *client.GH, prIssues []github.Issue, refresh bool) error {
	log.Println("creating a db connection")
	db, err := access.New()
	if err != nil {
		return err
	}
	users := user.NewRepo(db)
	if refresh {
		loadUsers(users, c)
		loadTeams(users, c)
	}
	events, err := access.NewEventAccess(db)
	if err != nil {
		return err
	}
	ch := make(chan event.Event, 10)
	var wg sync.WaitGroup
	go access.ReadAndPersist(events, ch, &wg)
	event.DumpEvents(c, prIssues, ch, &wg, users)
	wg.Wait()
	close(ch)
	return nil
}

func loadUsers(users *user.Repo, c *client.GH) {
	orgUsers, err := c.GetOrgUsers()
	if err != nil {
		log.Println(err)
	}
	for _, u := range orgUsers {
		err := users.SaveUser(u.GetLogin(), u.GetName())
		if err != nil {
			log.Println(err)
		}
	}
}

func loadTeams(users *user.Repo, c *client.GH) {
	for userId, teams := range c.GetTeamsByUser() {
		for _, team := range teams {
			err := users.SaveUserTeam(userId, team)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func printUsageAndExit() {
	log.Println("supply github organization and access token as command parameters")
	os.Exit(1)
}
