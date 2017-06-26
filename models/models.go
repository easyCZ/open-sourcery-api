package models

import (
	"github.com/google/go-github/github"
	"time"
)

type Issue struct {
	Id int
	Owner *github.User
	Repo string
	Title string
	Date *time.Time
	HtmlUrl string
}

func IssueFromGithubIssue(ghIssue *github.Issue) *Issue {
	return &Issue{
		Id: ghIssue.GetID(),
		Owner: ghIssue.Repository.Owner,
		Repo: *ghIssue.Repository.Name,
		Title: *ghIssue.Title,
		Date: ghIssue.CreatedAt,
		HtmlUrl: ghIssue.GetHTMLURL(),
	}
}