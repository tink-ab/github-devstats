# github-devstats

This application integrates with the GitHub API in order to collect information about merged Pull Requests within an organization.
Each PR that gets merged produces an event with metadata in the following format:
```go
e := Event{
	PrNumber:                 p.GetNumber(),
	Repository:               repo,
	MergedAt:                 p.GetMergedAt(),
	TimeToMergeSeconds:       p.GetMergedAt().Sub(p.GetCreatedAt()).Seconds(),
	BranchAgeSeconds:         branchAge(c, repo, p.GetBase().GetSHA(), p.GetMergeCommitSHA()),
	LinesAdded:               p.GetAdditions(),
	LinesRemoved:             p.GetDeletions(),
	FilesChanged:             p.GetChangedFiles(),
	CommitsCount:             p.GetCommits(),
	CommentsCount:            p.GetComments(),
	AuthorId:                 p.GetUser().GetLogin(),
	AuthorName:               c.GetUserName(p.GetUser().GetLogin()),
	AuthorTeams:              c.GetUserTeams(p.GetUser().GetLogin()),
	CommitsByType:            map[string]int{},
	FilesAddedByExtension:    map[string]int{},
	FilesModifiedByExtension: map[string]int{},
	TimeToApproveSeconds:     0,
	ApproverId:               "",
	ApproverName:             "",
	ApproverTeams:            nil,
	CrossTeam:                false,
	DismissReviewCount:       0,
	ChangesRequestedCount:    0,
}
```
which in turn can be forwarded to other time series systems such as Prometheus for further processing.
