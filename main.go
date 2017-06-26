package main

import (
	"github.com/opensourcery-io/api/updater"
	"flag"
	"github.com/opensourcery-io/api/services"
	"os"
	"github.com/golang/glog"
)

func verifyEnvVariables() {
	required := []string{
		services.ENV_GITHUB_CLIENT,
		services.ENV_GITHUB_SECRET,
	}

	for _, key := range required {
		val, found := os.LookupEnv(key)
		if !found {
			glog.Exitf("Environment variable %v is required. Found %v=%v", key, key, val)
		}
	}
}

func main() {
	verifyEnvVariables()

	fbCredsPath := flag.String("fbcreds", "", "The path to Firebase Credentials file")
	flag.Parse()

	if *fbCredsPath == "" {
		panic("Firebase Credentials are required, please see -h")
	}

	upd := updater.NewDefaultUpdater(*fbCredsPath)
	upd.Update()
}
