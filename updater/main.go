package updater

import (
	"fmt"
	"github.com/opensourcery-io/api/github"
)

func main() {
	service := github.NewGithubService()

	issues, _ := service.getIssuesWithLabels("facebook", "react-native", []string{"Help Wanted"})

	for _, issue := range issues {
		for _, label := range issue.Labels {
			fmt.Printf("%s\n", *label.Name)
		}
	}

}
