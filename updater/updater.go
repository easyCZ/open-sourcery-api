package updater

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/opensourcery-io/api/services"
	"github.com/google/go-github/github"
	"github.com/opensourcery-io/api/models"
	"strings"
	"github.com/golang/glog"
)

const (
	UPDATE_ACTION          = "update"
	PRINT_ACTION           = "print"
	DEFAULT_INDEX_FILEPATH = "./index.json"
	DEFAULT_ACTION         = UPDATE_ACTION
)

type ProjectsIndex map[string][]string

type Updater struct {
	Index           string // location of index file
	Action          string // print or write to firebase
	GithubService   *services.GithubService
	FirebaseService *services.FirebaseService
	LogosService    *services.LogoService
	SearchService   services.SearchService
}

func NewUpdater(
	index, action string,
	ghs *services.GithubService,
	lgs *services.LogoService,
	fbs *services.FirebaseService,
	ss services.SearchService) *Updater {

	return &Updater{
		Index:           index,
		Action:          action,
		GithubService:   ghs,
		LogosService:    lgs,
		FirebaseService: fbs,
		SearchService:   ss,
	}
}

func NewDefaultUpdater(firebaseCredsFile string) *Updater {
	fbs, err := services.NewFirebaseService(firebaseCredsFile)
	if err != nil {
		panic("Failed to create updater")
	}
	return NewUpdater(
		DEFAULT_INDEX_FILEPATH,
		DEFAULT_ACTION,
		services.NewGithubService(),
		services.NewLogoService(),
		fbs,
		services.NewAlgoliaSearch(),
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

func (u *Updater) updateProject(project string, labels []string, chn chan<- []*github.Issue) {
	glog.Infof("Processing %v", project)
	tokens := strings.SplitN(project, "/", 2)
	if len(tokens) < 2 {
		glog.Errorf("Failed to parse %v. Err: %v", project)
	}

	owner, repo := tokens[0], tokens[1]
	_, err := u.LogosService.Search(repo)
	if err != nil {
		glog.Warningf("Failed to get logo for %v, continuing. Err: %v", repo, err)
	}

	allIssues := make([]*github.Issue, 0)
	for _, label := range labels {
		issues, err := u.GithubService.GetIssuesWithLabels(owner, repo, []string{label})
		if err != nil {
			glog.Errorf("Failed to get allIssues for %v and label %v. Err: %v", project, label, err)
			continue
		}

		// transform to a models.Issue

		allIssues = append(allIssues, issues...)
	}
	chn <- allIssues
	glog.Infof("Finished %v", project)
}

func (u *Updater) Update() ([]*models.Issue, error) {
	index, err := u.loadProjectsIndex()
	if err != nil {
		glog.Errorf("Failed to parse projects index. Err: %v", err)
	}
	glog.Infof("Loaded projects index with %v items", len(*index))

	allIssuesChn := make(chan []*github.Issue, len(*index))
	defer close(allIssuesChn)

	for project, labels := range *index {
		go func(project string, labels []string) {
			u.updateProject(project, labels, allIssuesChn)
		}(project, labels)

	}

	allIssues := make([]*github.Issue, 0)
	for i := 0; i < len(*index); i++ {
		issues := <-allIssuesChn
		allIssues = append(allIssues, issues...)
	}

	limit, err := u.GithubService.GetRateLimit()
	glog.Infof("Github Rate Limit: %v", limit)

	errors := u.FirebaseService.StoreIssues(allIssues)
	errCount := 0
	for i := 0; i < len(allIssues); i++ {
		err := <-errors
		if err != nil {
			errCount += 1
		}
	}
	close(errors)

	if errCount > 0 {
		glog.Errorf("Encountered %v/%v errors when storing issues", errCount, len(allIssues))
	} else {
		glog.Infof("Stored %v issues", len(allIssues))
	}

	issuesToStore := make([]*models.Issue, 0)

	return issuesToStore, nil
}
