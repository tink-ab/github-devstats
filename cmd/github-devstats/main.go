package main

import (
	"fmt"
	"github.com/krlvi/github-devstats/client"
	"os"
)

func main() {
	org := os.Args[1]
	accessToken := os.Args[2]
	if len(org) <= 0 || len(accessToken) <= 0 {
		panic("supply github organization and access token as command parameters")
	}
	c := client.NewClient(org, accessToken)
	fmt.Println(c)
}
