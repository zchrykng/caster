package main

import (
	"fmt"
	"sort"
	"time"

	//"github.com/rogpeppe/godef/go/types"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/eduncan911/podcast"
	"github.com/gorilla/mux"
)

// Feed contains the information for each directory of episodes.
type Feed struct {
	URL      string
	Root     string
	Title    string
	Episodes map[string]*Episode
}

// EpisodeWrapper contains the file name and podcast.Item for each podcast file
type EpisodeWrapper struct {
	filename string
	item     podcast.Item
}

type byPosition []EpisodeWrapper

func (s byPosition) Len() int {
	return len(s)
}
func (s byPosition) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPosition) Less(i, j int) bool {
	return s[i].filename < s[j].filename
}

var typeFilter = regexp.MustCompile(".*\\.(m4a|m4b|mp3)")

// MakeFeed does the initial setup for the Feed structure and kicks off the initial scan for episodes
func MakeFeed(URL string, Root string, Title string) (*Feed, error) {
	f := &Feed{URL: URL, Root: Root, Title: Title}

	f.Episodes = make(map[string]*Episode)
	f.ScanEpisodes()

	return f, nil
}

// ScanEpisodes reads the root directory and creates an Episode for each matching media file in the directory
func (f *Feed) ScanEpisodes() error {
	files, err := ioutil.ReadDir(f.Root)
	if err != nil {
		log.Fatal(err)
	}

	for _, target := range files {
		if typeFilter.MatchString(target.Name()) {
			f.Episodes[Slugify(target.Name(), true)], err = MakeEpisode(path.Join(f.Root, target.Name()))
		}
	}

	return nil
}

// FeedHandler serves the RSS feed for podcast players
func (f *Feed) FeedHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	p := podcast.New(f.Title, f.URL, f.Title, &now, &now)

	items := []EpisodeWrapper{}

	for k, v := range f.Episodes {
		item := EpisodeWrapper{filename: v.Location}
		item.item = podcast.Item{
			Title:       v.Name,
			Link:        f.URL + "/" + k,
			Description: v.Name,
			PubDate:     &v.ModifiedTime,
		}

		item.item.AddEnclosure(item.item.Link, podcast.M4A, v.Size)

		items = append(items, item)
	}

	sort.Sort(byPosition(items))

	for _, v := range items {
		fmt.Println(v.item.Title)
		_, err := p.AddItem(v.item)
		if err != nil {
			fmt.Println(err)
		}
	}

	w.Header().Set("Content-Type", "application/xml")

	if err := p.Encode(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// FeedEpisode serves the media file when requested by a podcast player
func (f *Feed) FeedEpisode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["fileSlug"]

	http.ServeFile(w, r, f.Episodes[slug].Location)
}
