package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func basicAuth(pass handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			w.Header().Add("WWW-Authenticate", "Basic realm='rsvp'")
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			auth := strings.SplitN(header, " ", 2)

			if len(auth) != 2 || auth[0] != "Basic" {
				http.Error(w, "auhorization failed", http.StatusUnauthorized)
				return
			}

			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)

			if len(pair) != 2 || !validate(pair[0], pair[1]) {
				http.Error(w, "auhorization failed", http.StatusUnauthorized)
				return
			}

			pass(w, r)
		}
	}
}

func validate(username, password string) bool {
	if username == "Admin" && password == "Test123" {
		return true
	}
	return false
}
