package models

type RepositoryLabels struct {
	Repo   string `yaml: repo`
	Labels []string `yaml: labels`
}

type ProjectDef struct {
	Owner    string `yaml: owner`
	Projects []RepositoryLabels `yaml: projects`
}

type Issue struct {

}

