package updater

import (
	"github.com/opensourcery-io/api/models"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
	"github.com/golang/glog"
)

type ProjectsIndex map[string][]string

func LoadIndex(filepath string) ([]*models.OSProject, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	index := ProjectsIndex{}
	if err := yaml.Unmarshal(file, &index); err != nil {
		return nil, err
	}

	projects := make([]*models.OSProject, 0)
	for project, labels := range index {
		tokens := strings.SplitN(project, "/", 2)
		if len(tokens) < 2 {
			glog.Errorf("Failed to parse %v. Err: %v", project)
		}

		owner, repo := tokens[0], tokens[1]
		projects = append(projects, &models.OSProject{
			Owner: owner,
			Repo: repo,
			Labels: labels,
		})
	}
	glog.Infof("Loaded projects index with %v items", len(projects))

	return projects, nil
}
