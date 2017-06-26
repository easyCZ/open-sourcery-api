package services

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2/google"
	"gopkg.in/zabawaba99/firego.v1"
	"strconv"
	"github.com/golang/glog"
	"github.com/opensourcery-io/api/models"
)

const (
	FIREBASE_URL string = "https://open-sourcery.firebaseio.com/"
	ISSUES_KEY   string = "issues"
)

type FirebaseService struct {
	Service *firego.Firebase
}

func NewFirebaseService(credentialsPath string) (*FirebaseService, error) {

	credentials, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, errors.New("Failed to read firebase credentials")
	}

	conf, err := google.JWTConfigFromJSON(
		credentials,
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")

	return &FirebaseService{
		Service: firego.New(FIREBASE_URL, conf.Client(context.Background())),
	}, nil
}

func (fire *FirebaseService) getIssuesRef(key int) (*firego.Firebase, error) {
	return fire.Service.Ref(ISSUES_KEY + "/" + strconv.Itoa(key))
}

func (fire *FirebaseService) StoreIssue(issue *models.Issue) error {
	identifier := issue.HtmlUrl
	glog.Infof("Saving %v", identifier)

	issuesFb, err := fire.getIssuesRef(issue.Id)
	if err != nil {
		glog.Errorf("Failed to store %v. Err: %v", identifier, err)
		return err
	}

	glog.Infof("Saved %v.", identifier)
	return issuesFb.Set(issue)
}

func (fire *FirebaseService) StoreIssues(issues []*github.Issue) chan error {
	glog.Infof("Storing %v issues.", len(issues))
	collector := make(chan error, len(issues))

	for _, ghIssue := range issues {
		issue := models.IssueFromGithubIssue(ghIssue)
		go func(issue *models.Issue) {
			if issue == nil {
				glog.Error(issue)
			}

			err := fire.StoreIssue(issue)
			if err != nil {
				collector <- err
			}
			collector <- nil
		}(issue)
	}

	return collector
}
