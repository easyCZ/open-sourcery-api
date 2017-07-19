package updater

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/opensourcery-io/api/services"
	"github.com/opensourcery-io/api/models"
	"github.com/golang/glog"
	"fmt"
)

const (
	UPDATE_ACTION          = "update"
	PRINT_ACTION           = "print"
	DEFAULT_INDEX_FILEPATH = "./index.json"
	DEFAULT_ACTION         = UPDATE_ACTION
)

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

func (u *Updater) getIssuesForProject(owner, repo string, labels []string, chn chan<- []*models.Issue) {
	name := fmt.Sprintf("%v/%v", owner, repo)
	glog.Infof("Processing %v", name)

	_, err := u.LogosService.Search(repo)
	if err != nil {
		glog.Warningf("Failed to get logo for %v, continuing. Err: %v", repo, err)
	}

	allIssues := make([]*models.Issue, 0)
	for _, label := range labels {
		ghIssues, err := u.GithubService.GetIssuesWithLabels(owner, repo, []string{label})
		if err != nil {
			glog.Errorf("Failed to get allIssues for %v and label %v. Err: %v", name, label, err)
			continue
		}

		// transform to a models.Issue
		// may process duplicates but we don't care
		for _, ghIssue := range ghIssues {
			allIssues = append(allIssues, &models.Issue{
				Id:        ghIssue.GetID(),
				Owner:     owner,
				Repo:      repo,
				Title:     *ghIssue.Title,
				CreatedAt: ghIssue.CreatedAt,
				HtmlUrl:   ghIssue.GetHTMLURL(),
				Labels:    labels,
			})
		}
	}
	chn <- allIssues
	glog.Infof("Finished %v", name)
}

func (u *Updater) Update() ([]*models.Issue, error) {
	index, err := LoadIndex(u.Index)
	if err != nil {
		glog.Errorf("Failed to parse projects index. Err: %v", err)
		return nil, err
	}

	allIssuesChn := make(chan []*models.Issue, len(index))
	defer close(allIssuesChn)

	for _, project := range index {
		go func(owner, repo string, labels []string) {
			u.getIssuesForProject(owner, repo, labels, allIssuesChn)
		}(project.Owner, project.Repo, project.Labels)

	}

	// collect
	allIssues := make([]*models.Issue, 0)
	for i := 0; i < len(index); i++ {
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

	err = u.SearchService.IndexIssues(allIssues)
	if err != nil {
		glog.Errorf("Failed to index issues. Err: %v", err)
	}

	//err = u.FirebaseService.ReverseIndexIssuesByLabels(allIssues)
	//if err != nil {
	//	glog.Errorf("Failed to create reverse index. Err: %v", err)
	//}



	return allIssues, nil
}
