package services

import (
	"github.com/google/go-github/github"
	"os"
	"context"
	"github.com/golang/glog"
)

const (
	ENV_GITHUB_CLIENT string = "GITHUB_CLIENT"
	ENV_GITHUB_SECRET string = "GITHUB_SECRET"
)

type GithubService struct {
	client *github.Client
}

func NewGithubService() *GithubService {
	transport := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     os.Getenv(ENV_GITHUB_CLIENT),
		ClientSecret: os.Getenv(ENV_GITHUB_SECRET),
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

func (g *GithubService) GetRepo(owner, name string) (*github.Repository, error) {
	repo, _, err := g.client.Repositories.Get(context.Background(), owner, name)
	return repo, err
}

func (g *GithubService) GetRateLimit() (*github.RateLimits, error) {
	limits, _, err := g.client.RateLimits(context.Background())
	return limits, err
}
