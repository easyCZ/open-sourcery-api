package updater

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/opensourcery-io/api/services"
	"github.com/google/go-github/github"
	"github.com/opensourcery-io/api/models"
	"strings"
	"github.com/golang/glog"
	"fmt"
)

const (
	UPDATE_ACTION                 = "update"
	PRINT_ACTION                  = "print"
	DEFAULT_INDEX_FILEPATH        = "./index.json"
	DEFAULT_ACTION                = UPDATE_ACTION
	FIREBASE_CREDENTIALS_FILEPATH = "/Users/milanpavlik/golang/src/github.com/opensourcery-io/open-sourcery-firebase-adminsdk-f4ddp-83f1d4c231.json"
)

type ProjectsIndex map[string][]string

type Updater struct {
	Index           string // location of index file
	Action          string // print or write to firebase
	GithubService   *services.GithubService
	FirebaseService *services.FirebaseService
	LogosService    *services.LogosService
}

func NewUpdater(
	index, action string,
	ghs *services.GithubService,
	lgs *services.LogosService,
	fbs *services.FirebaseService) *Updater {

	return &Updater{
		Index:           index,
		Action:          action,
		GithubService:   ghs,
		LogosService:    lgs,
		FirebaseService: fbs,
	}
}

func NewDefaultUpdater() *Updater {
	//fbs, err := services.NewFirebaseService(FIREBASE_CREDENTIALS_FILEPATH)
	//if err != nil {
	//	panic("Failed to create updater")
	//}
	return NewUpdater(
		DEFAULT_INDEX_FILEPATH,
		DEFAULT_ACTION,
		services.NewGithubService(),
		services.NewLogosApiService(),
		nil,
	)
}

func (u *Updater) loadProjectsIndex() (*ProjectsIndex, error) {
	file, err := ioutil.ReadFile(u.Index)
	if err != nil {
		return nil, err
	}

	index := ProjectsIndex{}
	if err := yaml.Unmarshal(file, &index); err != nil {
		return nil, err
	}

	return &index, nil
}

func (u *Updater) Update() ([]*models.Issue, error) {
	index, err := u.loadProjectsIndex()
	allIssues := make([]*github.Issue, 0)

	for project, labels := range *index {
		tokens := strings.SplitN(project, "/", 2)
		if len(tokens) < 2 {
			glog.Errorf("Failed to parse %v. Err: %v", project, err)
			continue
		}

		owner, repo := tokens[0], tokens[1]
		//logo, err := u.LogosService.Search(repo)
		//if err != nil {
		//	glog.Warningf("Failed to get logo for %v, continuing. Err: %v", logo, err)
		//}

		for _, label := range labels {
			issues, err := u.GithubService.GetIssuesWithLabels(owner, repo, []string{label})
			if err != nil {
				glog.Errorf("Failed to get allIssues for %v and label %v. Err: %v", project, label, err)
				continue
			}

			allIssues = append(allIssues, issues...)
		}

	}

	limit, err := u.GithubService.GetRateLimit()
	fmt.Println(limit)

	//fmt.Println(allIssues)

	// TODO: Store to firebase
	issuesToStore := make([]*models.Issue, 0)

	return issuesToStore, nil
}
