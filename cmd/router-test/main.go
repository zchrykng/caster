package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	ttt := r.PathPrefix("/the-time-trilogy").Subrouter()

	ttt.HandleFunc("/", TimeHandler)

	http.Handle("/", r)

	http.ListenAndServe(":8000", nil)
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "The Time Trilogy")
}
