package caster

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Feed struct {
	URL      string
	Root     string
	Title    string
	Episodes map[string]*Episode
	Router   *mux.Router
}

func MakeFeed(URL string, Root string, Router *mux.Router, Title string) (*Feed, error) {
	f := &Feed{URL: URL, Root: Root, Title: Title, Router: Router}

	f.Episodes = make(map[string]*Episode)

	f.ScanEpisodes()

	return f, nil
}

func (f *Feed) ScanEpisodes() error {
	files, err := ioutil.ReadDir(f.Root)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println("Slug:", Slugify(f.Name(), true))
		fmt.Println("File:", f.Name())
	}

	return nil
}

func (f *Feed) FeedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Feed:", f.Title)
}

func (f *Feed) FeedEpisode(w http.ResponseWriter, r *http.Request) {

}
