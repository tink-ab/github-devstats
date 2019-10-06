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
