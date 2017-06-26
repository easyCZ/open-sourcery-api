package services

import (
	"github.com/opensourcery-io/api/models"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/fatih/structs"
	"os"
)

const (
	ENV_ALGOLIA_CLIENT = "ALGOLIA_CLIENT"
	ENV_ALGOLIA_SECRET = "ALGOLIA_SECRET"
)

type SearchService interface {
	IndexIssue(issue *models.Issue) error
}

type AlgoliaSearch struct {
	algoliasearch.Client
}

func (alg *AlgoliaSearch) getIssuesIndex() algoliasearch.Index {
	return alg.Client.InitIndex("issues")
}

func (alg *AlgoliaSearch) IndexIssue(issue *models.Issue) error {
	index := alg.getIssuesIndex()
	_, err := index.AddObject(structs.Map(issue))
	return err
}

func (alg *AlgoliaSearch) IndexIssues(issues []*models.Issue) error {
	converted := make([]algoliasearch.Object, 0, len(issues))
	for _, issue := range issues {
		converted = append(converted, structs.Map(issue))
	}
	index := alg.getIssuesIndex()
	_, err := index.AddObjects(converted)
	return err
}

func NewAlgoliaSearch() *AlgoliaSearch {
	return &AlgoliaSearch{
		algoliasearch.NewClient(os.Getenv(ENV_ALGOLIA_CLIENT), os.Getenv(ENV_ALGOLIA_SECRET)),
	}
}
