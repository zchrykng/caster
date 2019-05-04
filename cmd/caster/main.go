package main

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(c.Users[0])

	// r := mux.NewRouter()

	// http.Handle("/", r)

	// http.ListenAndServe(":8000", nil)
}
