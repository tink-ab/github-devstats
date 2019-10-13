package main

import (
	"github.com/krlvi/github-devstats/sql"
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

	toDB := false
	if len(os.Args) > 4 && os.Args[4] == "toDB" {
		toDB = true
	}

	log.Println("fetching github teams and their members for", org)
	c := client.NewClient(org, accessToken)

	date := time.Now().AddDate(0, 0, daysAgo*-1)
	log.Println("fetching merged pull requests for", org, "on date", date.Format("2006-01-02"))
	prIssues, err := c.GetAllMergedPRIssues(date)
	if err != nil {
		log.Panicln("could not fetch pull requests:", err)
	}

	if toDB {
		log.Println("writing", len(prIssues), "pull requests to db")
		db, err := sql.New()
		if err != nil {
			panic(err)
		}
		repo, err := sql.NewRepository(db)
		if err != nil {
			panic(err)
		}
		ch := make(chan event.Event, 10)
		var wg sync.WaitGroup
		go sql.ReadAndPersist(repo, ch, &wg)
		event.DumpEvents(c, prIssues, ch, &wg)
		wg.Wait()
		close(ch)
	} else {
		log.Println("outputting", len(prIssues), "pull requests to stdout")
		event.ProcessPRIssues(c, prIssues)
	}
}

func printUsageAndExit() {
	log.Println("supply github organization and access token as command parameters")
	os.Exit(1)
}
