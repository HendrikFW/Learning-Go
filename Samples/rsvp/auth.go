package main

import (
	"net/http"
)

func basicAuth(pass handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()

		if !ok {
			w.Header().Add("WWW-Authenticate", "Basic realm='rsvp'")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !validate(username, password) {
			http.Error(w, "auhorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func validate(username, password string) bool {
	return username == "Admin" && password == adminPw
}
