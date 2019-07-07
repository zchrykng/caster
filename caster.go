package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/gorilla/mux"
)

// Caster is the base struct for the podcast server.
type Caster struct {
	URL    string
	Root   string
	Feeds  map[string]*Feed
	Router *mux.Router
}

// MakeCaster takes a base url string and root directory string and returns a Caster struct
func MakeCaster(URL string, Root string) (*Caster, error) {
	c := &Caster{URL: URL, Root: Root}
	c.Feeds = make(map[string]*Feed)

	c.Router = mux.NewRouter()

	c.Router.HandleFunc("/", c.Handler)

	return c, nil
}

// ScanFeeds reads the files and folders contained in the root directory and creates Feed objects for each directory
func (c *Caster) ScanFeeds() error {
	files, err := ioutil.ReadDir(c.Root)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		slug := Slugify(f.Name(), true)
		if val, ok := c.Feeds[slug]; ok {
			val.ScanEpisodes()
		} else {
			if f.IsDir() {
				u := c.URL + "/" + slug

				c.Feeds[slug], err = MakeFeed(u, filepath.Join(c.Root, f.Name()), f.Name())
				if err != nil {
					log.Fatal(err)
				}

				c.Router.HandleFunc("/"+slug, c.Feeds[slug].FeedHandler)
				c.Router.HandleFunc("/"+slug+"/{fileSlug}", c.Feeds[slug].FeedEpisode)
			}
		}
	}

	fmt.Println("feeds:", c.Feeds)

	return nil
}

// Handler impliments the http.FuncHandler spec and serves a list of feeds
func (c *Caster) Handler(w http.ResponseWriter, r *http.Request) {
	// To store the keys in slice in sorted order
	var keys []string
	for k := range c.Feeds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// To perform the opertion you want
	for _, k := range keys {
		fmt.Fprintln(w, c.Feeds[k].Title, "<->", k)
	}
}
