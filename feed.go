package caster

import (
	"fmt"
	"html/template"
	//"github.com/rogpeppe/godef/go/types"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/gorilla/mux"
)

type Feed struct {
	URL      string
	Root     string
	Title    string
	Episodes map[string]*Episode
	Router   *mux.Router
}

var typeFilter = regexp.MustCompile(".*\\.(m4a|m4b|mp3)")

var feedTemplate = template.Must(template.ParseFiles("/User/zach/Documents/Projects/caster"))

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

	for _, target := range files {
		if typeFilter.MatchString(target.Name()) {
			f.Episodes[Slugify(target.Name(), true)], err = MakeEpisode(path.Join(f.Root, target.Name()))
		}
	}

	return nil
}

func (f *Feed) FeedHandler(w http.ResponseWriter, r *http.Request) {

	err := feedTemplate.Execute(w, f)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range f.Episodes {
		fmt.Fprintln(w, v.Name, "<->", k)
	}
}

func (f *Feed) FeedEpisode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["fileSlug"]

	http.ServeFile(w, r, f.Episodes[slug].Location)
}
