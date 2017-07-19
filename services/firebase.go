package services

import (
	"context"
	"errors"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"gopkg.in/zabawaba99/firego.v1"
	"strconv"
	"github.com/golang/glog"
	"github.com/opensourcery-io/api/models"
	"github.com/deckarep/golang-set"
	"strings"
)

const (
	FIREBASE_URL     string = "https://open-sourcery.firebaseio.com/"
	ISSUES_KEY       string = "issues"
	LABELS_TO_ISSUES string = "labelsToIssues"
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

func (fire *FirebaseService) getLabelsToIssuesRef() (*firego.Firebase, error) {
	return fire.Service.Ref(LABELS_TO_ISSUES)
}

func (fire *FirebaseService) StoreIssue(issue *models.Issue) error {
	identifier := issue.HtmlUrl
	glog.Infof("Saving %v", identifier)

	issuesFb, err := fire.getIssuesRef(issue.Id)
	if err != nil {
		glog.Errorf("Failed to get Issues ref %v. Err: %v", identifier, err)
		return err
	}

	glog.Infof("Saved %v.", identifier)
	return issuesFb.Set(issue)
}

func (fire *FirebaseService) ReverseIndexIssuesByLabels(issues []*models.Issue) error {
	glog.Infof("Reverse indexing %v issues", len(issues))

	index := make(map[string]mapset.Set, 0)
	for _, issue := range issues {
		id := issue.Id

		for _, label := range issue.Labels {

			normalized := strings.Replace(label, "$", " ", -1)
			normalized = strings.Replace(normalized, "#", " ", -1)
			normalized = strings.Replace(normalized, "[", " ", -1)
			normalized = strings.Replace(normalized, "]", " ", -1)
			normalized = strings.Replace(normalized, "/", " ", -1)
			normalized = strings.Replace(normalized, ".", " ", -1)
			normalized = strings.TrimSpace(normalized)
			//glog.Infof("Storing %v and label %v", id, label)
			if _, ok := index[normalized]; !ok {
				index[normalized] = mapset.NewSet()
			}
			index[normalized].Add(id)
		}
	}

	fb, err := fire.getLabelsToIssuesRef()
	if err != nil {
		glog.Errorf("Failed to get reverse index for labels ref. Err: %v", err)
		return err
	}
	return fb.Set(index)
}

func (fire *FirebaseService) StoreIssues(issues []*models.Issue) chan error {
	glog.Infof("Storing %v issues.", len(issues))
	collector := make(chan error, len(issues))

	for _, issue := range issues {
		go func(issue *models.Issue) {
			if issue == nil {
				glog.Error(issue)
			}

			err := fire.StoreIssue(issue)
			if err != nil {
				collector <- err
			}

			// store reverse index

			collector <- nil
		}(issue)
	}

	return collector
}
