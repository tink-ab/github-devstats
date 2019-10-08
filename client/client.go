package client

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GH struct {
	client *github.Client
	ctx    context.Context
	org    string
	users  map[string]string
	teams  map[string][]string
}

type rateLimitRetryer struct {
	Delegate http.RoundTripper
}

const MaxRetries = 10

func (r rateLimitRetryer) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	retries := 0
	for {
		if retries > MaxRetries {
			return nil, errors.New("too many retries")
		}

		resp, err = r.Delegate.RoundTrip(req)
		if err != nil {
			return
		}
		if resp.StatusCode != 403 {
			return
		}
		if resp.Header.Get("X-RateLimit-Remaining") != "0" {
			// some other error, let's not back-off in that case.
			return
		}

		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		// From Github: Mon, 01 Jul 2013 17:27:06 GMT
		// From Golang: Mon Jan 2 15:04:05 -0700 MST 2006
		now, parseErr := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", resp.Header.Get("Date"))
		if parseErr != nil {
			return nil, parseErr
		}

		resetTime, resetTimeErr := strconv.Atoi(resp.Header.Get("X-RateLimit-Reset"))
		if resetTimeErr != nil {
			return nil, resetTimeErr
		}
		reset := time.Unix(int64(resetTime), 0)

		timeToSleep := reset.Sub(now)
		if timeToSleep < 0 {
			// Not sure why this would happen but better safe than sorry.
			timeToSleep = 0
		}
		timeToSleep += time.Second
		log.Println("Sleeping due to rate-limitting:", timeToSleep)
		time.Sleep(timeToSleep)

		retries++
	}
}

func NewClient(org, accessToken string) *GH {
	ctx := context.Background()
	httpClient := &http.Client{
		Transport: rateLimitRetryer{
			&oauth2.Transport{
				Base: http.DefaultTransport,
				Source: oauth2.ReuseTokenSource(nil, oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: accessToken},
				)),
			},
		},
	}

	client := &GH{
		client: github.NewClient(httpClient),
		ctx:    ctx,
		org:    org,
		users:  map[string]string{},
	}
	client.teams = client.GetTeamsByUser()
	return client
}

func (c *GH) GetUserTeams(login string) []string {
	return c.teams[login]
}

func (c *GH) GetTeamsByUser() map[string][]string {
	teams, err := c.GetTeams()
	if err != nil {
		return nil
	}
	membership := map[string]map[string]bool{}
	for _, t := range teams {
		members, err := c.GetMembers(t.GetID())
		if err != nil {
			continue
		}
		for _, m := range members {
			if _, ok := membership[m.GetLogin()]; !ok {
				membership[m.GetLogin()] = map[string]bool{}
			}
			membership[m.GetLogin()][t.GetSlug()] = true
		}
	}
	out := map[string][]string{}
	for user, teams := range membership {
		out[user] = []string{}
		for t := range teams {
			out[user] = append(out[user], t)
		}
	}
	return out
}

func (c *GH) GetTeams() ([]*github.Team, error) {
	var teams []*github.Team
	page := 1
	for page != 0 {
		newTeams, nextPage, err := getTeams(c, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		teams = append(teams, newTeams...)
	}
	return teams, nil
}

func getTeams(c *GH, page int) (teams []*github.Team, nextPage int, err error) {
	teams, rsp, err := c.client.Teams.ListTeams(c.ctx, c.org, &github.ListOptions{Page: page})
	if err != nil {
		return nil, 0, err
	}
	return teams, rsp.NextPage, nil
}

func (c *GH) GetMembers(teamId int64) ([]*github.User, error) {
	var members []*github.User
	page := 1
	for page != 0 {
		newMembers, nextPage, err := getMembers(c, teamId, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		members = append(members, newMembers...)
	}
	return members, nil
}

func getMembers(c *GH, teamId int64, page int) (members []*github.User, nextPage int, err error) {
	members, rsp, err := c.client.Teams.ListTeamMembers(c.ctx, teamId, &github.TeamListTeamMembersOptions{
		ListOptions: github.ListOptions{Page: page}})
	if err != nil {
		return nil, 0, err
	}
	return members, rsp.NextPage, nil
}

func (c *GH) GetCommit(repo, hash string) (*github.Commit, error) {
	commit, _, err := c.client.Git.GetCommit(c.ctx, c.org, repo, hash)
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func (c *GH) GetPRCommits(prNumber int, repo string) ([]*github.RepositoryCommit, error) {
	var commits []*github.RepositoryCommit
	page := 1
	for page != 0 {
		newCommits, nextPage, err := getCommits(c, prNumber, repo, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		commits = append(commits, newCommits...)
	}
	return commits, nil
}

func getCommits(c *GH, prNumber int, repo string, page int) (commits []*github.RepositoryCommit, nextPage int, err error) {
	commits, rsp, err := c.client.PullRequests.ListCommits(c.ctx, c.org, repo, prNumber, &github.ListOptions{Page: page})
	if err != nil {
		return nil, 0, err
	}
	return commits, rsp.NextPage, nil
}

func (c *GH) GetPRFiles(prNumber int, repo string) ([]*github.CommitFile, error) {
	var files []*github.CommitFile
	page := 1
	for page != 0 {
		newFiles, nextPage, err := getFiles(c, prNumber, repo, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		files = append(files, newFiles...)
	}
	return files, nil
}

func getFiles(c *GH, prNumber int, repo string, page int) (files []*github.CommitFile, nextPage int, err error) {
	files, rsp, err := c.client.PullRequests.ListFiles(c.ctx, c.org, repo, prNumber, &github.ListOptions{Page: page})
	if err != nil {
		return nil, 0, err
	}
	return files, rsp.NextPage, nil
}

func (c *GH) GetReviews(prNumber int, repo string) ([]*github.PullRequestReview, error) {
	var reviews []*github.PullRequestReview
	page := 1
	for page != 0 {
		newReviews, nextPage, err := getReviews(c, prNumber, repo, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		reviews = append(reviews, newReviews...)
	}
	return reviews, nil
}

func getReviews(c *GH, prNumber int, repo string, page int) (review []*github.PullRequestReview, nextPage int, err error) {
	reviews, rsp, err := c.client.PullRequests.ListReviews(c.ctx, c.org, repo, prNumber, &github.ListOptions{Page: page})
	if err != nil {
		return nil, 0, err
	}
	return reviews, rsp.NextPage, nil
}

func (c *GH) GetUserName(user string) string {
	if len(c.users[user]) > 0 {
		return c.users[user]
	}
	u, _, err := c.client.Users.Get(c.ctx, user)
	if err != nil {
		return ""
	}
	c.users[user] = u.GetName()
	return u.GetName()
}

func (c *GH) GetPR(prNumber int, repo string) (*github.PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(c.ctx, c.org, repo, prNumber)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (c *GH) GetMergedPRs(date time.Time) (prs []*github.PullRequest, reposByPR map[int]string, err error) {
	issues, err := getAllMergedPRIssues(c, date)
	if err != nil {
		return nil, nil, err
	}
	reposByPR = map[int]string{}
	for _, i := range issues {
		repo := repoUrlToName(i.GetRepositoryURL())
		reposByPR[i.GetNumber()] = repo
		pr, err := c.GetPR(i.GetNumber(), repo)
		if err != nil {
			continue
		}
		prs = append(prs, pr)
	}
	return prs, reposByPR, nil
}

func getAllMergedPRIssues(c *GH, date time.Time) ([]github.Issue, error) {
	var prs []github.Issue
	page := 1
	for page != 0 {
		newPrs, nextPage, err := getMergedPRIssues(c, date, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		prs = append(prs, newPrs...)
	}
	return prs, nil
}

func getMergedPRIssues(c *GH, date time.Time, page int) (prs []github.Issue, nextPage int, err error) {
	d := date.Format("2006-01-02")
	query := "org:" + c.org + " is:pr is:closed merged:" + d + ".." + d
	issues, rsp, err := c.client.Search.Issues(c.ctx, query, &github.SearchOptions{ListOptions: github.ListOptions{Page: page}})
	if err != nil {
		return nil, 0, err
	}
	return issues.Issues, rsp.NextPage, nil
}

func repoUrlToName(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}
