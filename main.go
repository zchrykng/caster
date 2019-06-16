package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mitchellh/go-homedir"
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

	hostPtr := flag.String("host", "localhost", "host name")
	portPtr := flag.String("port", "8000", "port")
	rootPtr := flag.String("root", "~/.config/caster/casts", "cast root folder")
	userfilePtr := flag.String("userfile", "~/.config/caster/users.json", "user config file")

	flag.Parse()

	c := Config{}
	c.URL = fmt.Sprintf("%s:%s", *hostPtr, *portPtr)
	c.Root, _ = homedir.Expand(*rootPtr)

	userPath, _ := homedir.Expand(*userfilePtr)
	userfile, err := os.Open(userPath)
	defer userfile.Close()
	if err != nil {
		fmt.Println("error:", err)
	}
	userDecoder := json.NewDecoder(userfile)

	err = userDecoder.Decode(&c.Users)
	if err != nil {
		fmt.Println("error:", err)
	}

	cast, err := MakeCaster(c.URL, c.Root)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = cast.ScanFeeds()
	if err != nil {
		fmt.Println("error:", err)
	}

	http.Handle("/", cast.Router)

	http.ListenAndServe(fmt.Sprintf(":%s", *portPtr), nil)
}
