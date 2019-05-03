package gocaster

import "time"

type Feed struct {
	Title       string
	Link        string
	Description string
	Episodes    []Episode
}

type Episode struct {
	Title string
	File  File
}

type File struct {
	URL          string
	Name         string
	Artist       string
	ModifiedTime time.Time
	Duration     time.Duration
	Type         string
	Size         int
}
