package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
)

type BasicAuth struct {
	users map[string]string
}

func MakeBasicAuth(users []*User) *BasicAuth {
	ba := &BasicAuth{}

	ba.users = make(map[string]string)

	err := ba.Populate(users)
	if err != nil {
		log.Fatal(err)
	}

	return ba
}

func (ba *BasicAuth) Populate(users []*User) error {

	for _, v := range users {
		ba.users[v.Name] = v.Pass
	}

	return nil
}

func (ba *BasicAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(pass), []byte(ba.users[user])) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Podcast Login"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorised.\n"))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

