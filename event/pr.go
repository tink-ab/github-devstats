package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/krlvi/github-devstats/client"
)

type Event struct {
	PrNumber                 int            `json:"pr_number"`
	Repository               string         `json:"repository"`
	MergedAt                 time.Time      `json:"merged_at"`
	TimeToMerge              Duration       `json:"time_to_merge"`
	LinesAdded               int            `json:"lines_added"`
	LinesRemoved             int            `json:"lines_removed"`
	FilesChanged             int            `json:"files_changed"`
	CommitsCount             int            `json:"commits_count"`
	CommentsCount            int            `json:"comments_count"`
	AuthorId                 string         `json:"author_id"`
	AuthorName               string         `json:"author_name"`
	AuthorTeams              []string       `json:"author_teams"`
	CommitsByType            map[string]int `json:"commits_by_type"`
	FilesAddedByExtension    map[string]int `json:"files_added_by_extension"`
	FilesModifiedByExtension map[string]int `json:"files_modified_by_extension"`
	TimeToApprove            Duration       `json:"time_to_approve"`
	ApproverId               string         `json:"approver_id"`
	ApproverName             string         `json:"approver_name"`
	ApproverTeams            []string       `json:"approver_teams"`
	CrossTeam                bool           `json:"cross_team"`
	DismissReviewCount       int            `json:"dismiss_review_count"`
	ChangesRequestedCount    int            `json:"changes_requested_count"`
}

func ProcessPRs(c *client.GH, prs []*github.PullRequest, prRepos map[int]string) {
	for _, p := range prs {
		j, err := json.MarshalIndent(prToEvent(c, p, prRepos), "", "  ")
		if err != nil {
			continue
		}
		fmt.Printf("%s\n", string(j))
	}
}

func prToEvent(c *client.GH, p *github.PullRequest, prRepos map[int]string) Event {
	e := Event{
		PrNumber:                 p.GetNumber(),
		Repository:               prRepos[p.GetNumber()],
		MergedAt:                 p.GetMergedAt(),
		TimeToMerge:              Duration{p.GetMergedAt().Sub(p.GetCreatedAt())},
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
		TimeToApprove:            Duration{},
		ApproverId:               "",
		ApproverName:             "",
		ApproverTeams:            nil,
		CrossTeam:                false,
		DismissReviewCount:       0,
		ChangesRequestedCount:    0,
	}

	commits, err := c.GetPRCommits(p.GetNumber(), prRepos[p.GetNumber()])
	if err == nil {
		for _, com := range commits {
			e.CommitsByType[commitType(com.GetCommit().GetMessage())]++
		}
	}

	files, err := c.GetPRFiles(p.GetNumber(), prRepos[p.GetNumber()])
	if err == nil {
		for _, f := range files {
			if f.GetStatus() == "modified" {
				e.FilesModifiedByExtension[fileExtension(f.GetFilename())]++
			}
			if f.GetStatus() == "added" {
				e.FilesAddedByExtension[fileExtension(f.GetFilename())]++
			}
		}
	}

	reviews, err := c.GetReviews(p.GetNumber(), prRepos[p.GetNumber()])
	if err == nil {
		for _, r := range reviews {
			if r.GetState() == "APPROVED" {
				e.TimeToApprove = Duration{r.GetSubmittedAt().Sub(p.GetCreatedAt())}
				e.ApproverId = r.GetUser().GetLogin()
				e.ApproverName = c.GetUserName(r.GetUser().GetLogin())
				e.ApproverTeams = c.GetUserTeams(r.GetUser().GetLogin())
				e.CrossTeam = crossTeam(c.GetUserTeams(p.GetUser().GetLogin()), c.GetUserTeams(r.GetUser().GetLogin()))
			}
			if r.GetState() == "DISMISSED" {
				e.DismissReviewCount++
			}
			if r.GetState() == "CHANGES_REQUESTED" {
				e.ChangesRequestedCount++
			}
		}
	}
	return e
}

func crossTeam(from, to []string) bool {
	fromSet := map[string]bool{}
	for _, f := range from {
		fromSet[f] = true
	}
	for _, t := range to {
		if fromSet[t] {
			return false
		}
	}
	return true
}

func fileExtension(filename string) string {
	tokens := strings.FieldsFunc(filename, delimiter)
	return tokens[len(tokens)-1]
}

func delimiter(r rune) bool {
	return r == '.' || r == '/'
}

func commitType(msg string) string {
	if strings.HasPrefix(msg, "build") {
		return "build"
	}
	if strings.HasPrefix(msg, "chore") {
		return "chore"
	}
	if strings.HasPrefix(msg, "ci") {
		return "ci"
	}
	if strings.HasPrefix(msg, "copy") {
		return "copy"
	}
	if strings.HasPrefix(msg, "doc") {
		return "docs"
	}
	if strings.HasPrefix(msg, "fea") {
		return "feat"
	}
	if strings.HasPrefix(msg, "fix") {
		return "fix"
	}
	if strings.HasPrefix(msg, "log") {
		return "log"
	}
	if strings.HasPrefix(msg, "perf") {
		return "perf"
	}
	if strings.HasPrefix(msg, "ref") {
		return "refactor"
	}
	if strings.HasPrefix(msg, "revert") {
		return "revert"
	}
	if strings.HasPrefix(msg, "style") {
		return "style"
	}
	if strings.HasPrefix(msg, "test") {
		return "test"
	}
	return "uncategorized"
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
