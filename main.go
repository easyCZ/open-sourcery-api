package main

import (
	"github.com/opensourcery-io/api/updater"
	"flag"
)

func main() {
	fbCredsPath := flag.String("fbcreds", "", "The path to Firebase Credentials file")
	flag.Parse()

	if *fbCredsPath == "" {
		panic("Firebase Credentials are required, please see -h")
	}

	upd := updater.NewDefaultUpdater(*fbCredsPath)
	upd.Update()
}
