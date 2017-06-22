package main

import (
	"fmt"
	"io/ioutil"
	"github.com/opensourcery-io/api/models"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"strings"
	"github.com/opensourcery-io/api/services"
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



func transform() {
	jsonFile := "/Users/milanpavlik/golang/src/github.com/opensourcery-io/api/labels_by_repos_relevant_2.json"
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	// { label: ["owner/repo"], ... }
	var parsed map[string][]string

	if err := json.Unmarshal(data, &parsed); err != nil {
		panic(err)
	}

	for label, repos := range parsed {
		for _, repoPath := range repos {
			tokens := strings.SplitN(repoPath, "/", 1)
			owner := tokens[0]
			repository := tokens[1]

			def, err := readProjectDef("/Users/milanpavlik/golang/src/github.com/opensourcery-io/api/projects" + owner + ".yml")
			if err != nil {
				def = &models.ProjectDef{
					Owner: owner,
					Projects: []models.RepositoryLabels{},
				}
			}


			repoLabel := models.RepositoryLabels{
				Repo: repository,

			}
			def.Projects = append(def.Projects, )

		}
	}

	gh := services.NewGithubService()

	for label, repos := range parsed {

		projects := make([]models.RepositoryLabels, 0)

		for _, name := range repos {


			repo, err := gh.GetRepo(owner, repository)
			if err != nil {
				panic(err)
			}
			project := models.RepositoryLabels{
				Repo: *repo.Name,

			}



		}

		def := &models.ProjectDef{
			Owner: *repo.Owner.Name,
			Projects: []models.RepositoryLabels{
				models.RepositoryLabels{
					Repo:
				},
			},
		}
	}

	fmt.Println(parsed)

}

func main() {
	//fb, err := services.NewFirebaseService("/Users/milanpavlik/golang/src/github.com/opensourcery-io/open-sourcery-firebase-adminsdk-f4ddp-83f1d4c231.json")

	//gh := services.NewGithubService()
	//
	//dir, err := os.Getwd()
	//if err != nil {
	//	fmt.Errorf("Failed to get cwd %v", err)
	//}
	//projectDefs, err := listProjectDefs(dir + "/projects")
	//issues := make([]*github.Issue, 0)
	//
	//for _, projectDef := range projectDefs {
	//	for _, project := range projectDef.Projects {
	//		for _, label := range project.Labels {
	//			iss, _ := gh.GetIssuesWithLabels(projectDef.Owner, project.Repo, []string{label})
	//			issues = append(issues, iss...)
	//		}
	//
	//	}
	//}
	//
	//for _, issue := range issues {
	//	if err := fb.StoreIssue(issue); err != nil {
	//		fmt.Errorf("Failed to store issue %v", err)
	//	}
	//
	//}
	transform()

	//logosService := services.NewLogosApiService()
	//logo, _ := logosService.Search("facebook")
	//
	//fmt.Printf("%v", logo)

	//

	//
	//for _, issue := range issues {
	//	for _, label := range issue.Labels {
	//		fmt.Printf("%s\n", *label.Name)
	//	}
	//}

}
