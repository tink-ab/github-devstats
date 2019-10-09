# github-devstats

This application integrates with the GitHub API in order to collect information about merged Pull Requests within an organization.
Each PR that gets merged produces an event with metadata in the following format:
```go
type Event struct {
	PrNumber                 int            `json:"pr_number"`
	Repository               string         `json:"repository"`
	MergedAt                 time.Time      `json:"merged_at"`
	TimeToMergeSeconds       float64        `json:"time_to_merge_seconds"`
	BranchAgeSeconds         float64        `json:"branch_age_seconds"`
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
	JavaTestFilesModified    int            `json:"java_test_files_modified"`
	JavaTestsAdded           int            `json:"java_tests_added"`
	TimeToApproveSeconds     float64        `json:"time_to_approve_seconds"`
	ApproverId               string         `json:"approver_id"`
	ApproverName             string         `json:"approver_name"`
	ApproverTeams            []string       `json:"approver_teams"`
	CrossTeam                bool           `json:"cross_team"`
	DismissReviewCount       int            `json:"dismiss_review_count"`
	ChangesRequestedCount    int            `json:"changes_requested_count"`
}
```
which in turn can be forwarded to other time series systems such as Prometheus for further processing.

## Prerequisites
Have [Bazel](https://bazel.build/) installed, ideally with [Bazelisk](https://github.com/bazelbuild/bazelisk) so that it picks up the version from the [.bazelversion](.bazelversion) file automatically.

## Usage
Run:

`$ bazel run //cmd/github-devstats your-gh-org $GH_TOKEN`

Where `your-gh-org` is the name of your GitHub org and `$GH_TOKEN` evaluates to your secret token (you probably dont want to put this in your shell history).
