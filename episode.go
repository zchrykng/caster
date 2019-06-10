package caster

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dhowden/tag"
)

// Episode contains the information about the podcast episode
type Episode struct {
	Slug         string
	Location     string
	Name         string
	Artist       string
	ModifiedTime time.Time
	Duration     time.Duration
	Type         tag.FileType
	Size         int64
}

// MakeEpisode takes a file location and returns an episode object based on the meta data
func MakeEpisode(location string) (*Episode, error) {
	e := &Episode{}

	f, err := os.Open(location)
	if err != nil {
		fmt.Printf("Episode handler: error loading file: %v", err)
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("Episode handler: error reading file info: %v", err)
		return nil, err
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		log.Fatal(err)
	}

	e.Slug = Slugify(m.Title(), true)
	e.Location = location
	e.Name = m.Title()
	e.Artist = m.Artist()
	e.ModifiedTime = fi.ModTime()
	e.Type = m.FileType()
	e.Size = fi.Size()

	return e, nil
}
