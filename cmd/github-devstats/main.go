package main

import (
	"github.com/krlvi/github-devstats/client"
	"github.com/krlvi/github-devstats/event"
	"os"
	"time"
)

func main() {
	org := os.Args[1]
	accessToken := os.Args[2]
	if len(org) <= 0 || len(accessToken) <= 0 {
		panic("supply github organization and access token as command parameters")
	}
	c := client.NewClient(org, accessToken)

	yesterday := time.Now().AddDate(0, 0, -1)
	prs, reposByPR, err := c.GetMergedPRs(yesterday)
	if err != nil {
		panic("could not fetch pull requests")
	}
	event.ProcessPRs(c, prs, reposByPR)
}
