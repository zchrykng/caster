package caster

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Caster struct {
	URL    string
	Root   string
	Feeds  map[string]*Feed
	Router *mux.Router
}

func MakeCaster(URL string, Root string) (*Caster, error) {
	c := &Caster{URL: URL, Root: Root}
	c.Feeds = make(map[string]*Feed)

	c.Router = mux.NewRouter()

	c.Router.HandleFunc("/", c.Handler)

	return c, nil
}

func (c *Caster) ScanFeeds() error {
	files, err := ioutil.ReadDir(c.Root)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		slug := Slugify(f.Name(), true)
		u := path.Join(c.URL, slug)

		s := c.Router.PathPrefix("/" + slug).Subrouter()

		c.Feeds[slug], err = MakeFeed(u, filepath.Join(c.Root, f.Name()), s, f.Name())
		if err != nil {
			log.Fatal(err)
		}

		c.Router.HandleFunc("/"+slug, c.Feeds[slug].FeedHandler)
		c.Router.HandleFunc("/"+slug+"/{fileSlug}", c.Feeds[slug].FeedEpisode)
	}

	fmt.Println("feeds:", c.Feeds)

	return nil
}

func (c *Caster) Handler(w http.ResponseWriter, r *http.Request) {

	for k, v := range c.Feeds {
		fmt.Fprintln(w, v.Title, "<->", k)
	}
}
