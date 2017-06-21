package github

import (
	"github.com/google/go-github/github"
	"os"
	"context"
	"github.com/golang/glog"
)

type GithubService struct {
	client *github.Client
}

func NewGithubService() *GithubService {
	transport := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     os.Getenv("GITHUB_CLIENT"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
	}
	return &GithubService{
		client: github.NewClient(transport.Client()),
	}
}

func (g *GithubService) GetIssuesWithLabels(owner, repository string, labels []string) ([]*github.Issue, error) {
	query := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
		Labels: labels,
	}
	issues, resp, err := g.client.Issues.ListByRepo(context.Background(), owner, repository, query)
	if err != nil {
		glog.Error("Failed to get issues with labels", err, resp)
	}
	return issues, err
}
