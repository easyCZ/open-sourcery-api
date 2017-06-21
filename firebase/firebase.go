package firebase

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2/google"
	"gopkg.in/zabawaba99/firego.v1"
	"strconv"
)

const FIREBASE_URL string = "https://open-sourcery.firebaseio.com/"

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
	return fire.Service.Ref("issues2/" + strconv.Itoa(key))
}

func (fire *FirebaseService) StoreIssue(issue *github.Issue) error {
	issuesFb, err := fire.getIssuesRef(issue.GetID())
	if err != nil {
		return err
	}
	return issuesFb.Set(issue)
}
