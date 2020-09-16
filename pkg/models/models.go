package models

import "time"

type Story struct {
	ID      int
	Title   string
	Content string
	Created time.Time
}

type Stories []*Story
