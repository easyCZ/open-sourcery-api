package models

import (
	"time"
)

type Issue struct {
	Id        int
	Owner     string
	Repo      string
	Title     string
	CreatedAt *time.Time
	HtmlUrl   string
	Labels    []string
}

type OSProject struct {
	Owner  string
	Repo   string
	Labels []string
}
