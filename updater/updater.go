package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/opensourcery-io/api/models"
	"gopkg.in/yaml.v2"
	"github.com/opensourcery-io/api/services"
	"github.com/google/go-github/github"
)

func readProjectDef(filepath string) (*models.ProjectDef, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	p := models.ProjectDef{}

	if err := yaml.Unmarshal(file, &p); err != nil {
		return nil, err
	}

	fmt.Print(p)

	return &p, nil
}

func listProjectDefs(dir string) ([]*models.ProjectDef, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	defs := make([]*models.ProjectDef, 0)
	for _, file := range files {
		if !file.IsDir() {
			filename := dir + "/" + file.Name()
			projectDef, err := readProjectDef(filename)
			if err != nil {
				return nil, err
			}
			defs = append(defs, projectDef)
		}

	}

	return defs, nil
}

func main() {
	fb, err := services.NewFirebaseService("/Users/milanpavlik/golang/src/github.com/opensourcery-io/open-sourcery-firebase-adminsdk-f4ddp-83f1d4c231.json")

	gh := services.NewGithubService()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Failed to get cwd %v", err)
	}
	projectDefs, err := listProjectDefs(dir + "/projects")
	issues := make([]*github.Issue, 0)

	for _, projectDef := range projectDefs {
		for _, project := range projectDef.Projects {
			for _, label := range project.Labels {
				iss, _ := gh.GetIssuesWithLabels(projectDef.Owner, project.Repo, []string{label})
				issues = append(issues, iss...)
			}

		}
	}

	for _, issue := range issues {
		if err := fb.StoreIssue(issue); err != nil {
			fmt.Errorf("Failed to store issue %v", err)
		}

	}

	logosService := services.NewLogosApiService()
	logo, _ := logosService.Search("facebook")

	fmt.Printf("%v", logo)

	for _, issue := range issues {
		for _, label := range issue.Labels {
			fmt.Printf("%s\n", *label.Name)
		}
	}

}
