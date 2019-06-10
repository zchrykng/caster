package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mitchellh/go-homedir"
	caster "github.com/zchrykng/gocaster"
)

type User struct {
	Name        string
	Pass        string
	Admin       bool
	AuthedFeeds []string
}

type Config struct {
	URL   string
	Root  string
	Users []*User
}

func main() {

	path, _ := homedir.Expand("~/.caster")

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println("error:", err)
	}
	decoder := json.NewDecoder(file)
	c := Config{}
	err = decoder.Decode(&c)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(c)

	c.Root, _ = homedir.Expand(c.Root)

	cast, err := caster.MakeCaster(c.URL, c.Root)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = cast.ScanFeeds()
	if err != nil {
		fmt.Println("error:", err)
	}

	http.Handle("/", cast.Router)

	http.ListenAndServe(":8000", nil)
}
