package main

import "github.com/opensourcery-io/api/updater"

func main() {
	upd := updater.NewDefaultUpdater()

	upd.Update()
}
