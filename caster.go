package caster

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
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

	return c, nil
}

func (c *Caster) ScanFeeds() error {

	path, _ := homedir.Expand(c.Root)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		u := c.URL + "/" + Slugify(f.Name(), true)
		c.Feeds[Slugify(f.Name(), true)], err = MakeFeed(u, c.Root, f.Name())

		fmt.Println(f.Name())
	}

	fmt.Println("feeds:", c.Feeds)

	return nil
}

func (c *Caster) Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	feed := vars["feed"]

	fmt.Println("episodes:", c.Feeds[feed])
}
