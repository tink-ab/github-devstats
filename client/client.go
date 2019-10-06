package client

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GH struct {
	client *github.Client
	ctx    context.Context
	org    string
}

func NewClient(org, accessToken string) *GH {
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	))
	client := &GH{
		client: github.NewClient(httpClient),
		ctx:    ctx,
		org:    org,
	}
	return client
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
