package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

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
	userFilePtr := flag.String("userfile", "~/.config/caster/users.json", "user config file")
	exposePortPtr := flag.Bool("exposeport", false, "if the port should be included in the host url")

	flag.Parse()

	c := Config{}
	if *exposePortPtr {
		c.URL = fmt.Sprintf("https://%s:%s", *hostPtr, *portPtr)
	} else {
		c.URL = fmt.Sprintf("https://%s", *hostPtr)
	}

	c.Root, _ = homedir.Expand(*rootPtr)

	userPath, _ := homedir.Expand(*userFilePtr)
	userFile, err := os.Open(userPath)
	defer userFile.Close()
	if err != nil {
		fmt.Println("error:", err)
	}
	userDecoder := json.NewDecoder(userFile)

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

	ticker := time.NewTicker(15 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := cast.ScanFeeds()
				if err != nil {
					fmt.Println("error:", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	ba := MakeBasicAuth(c.Users)

	cast.Router.Use(ba.Middleware)

	http.Handle("/", cast.Router)

	http.ListenAndServe(fmt.Sprintf(":%s", *portPtr), nil)
}
