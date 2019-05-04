package caster

import "time"

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
