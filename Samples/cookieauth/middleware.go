package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type user struct {
	Authenticated bool
	UserName      string
}

type handler func(w http.ResponseWriter, r *http.Request)

type authHandler func(u user, w http.ResponseWriter, r *http.Request)

func method(next handler, allowed ...string) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if contains(allowed, r.Method) {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func secure(next authHandler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("signin")
		if err != nil {
			redirectURL := fmt.Sprintf("/signin?returnURL=%s", url.QueryEscape(r.URL.Path))
			w.Header().Add("Location", redirectURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			u := user{
				UserName: cookie.Value,
			}
			next(u, w, r)
		}
	}
}

func allowAnonymous(next authHandler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("signin")
		if err != nil {
			u := user{
				Authenticated: false,
			}
			next(u, w, r)
		} else {
			u := user{
				Authenticated: true,
				UserName:      cookie.Value,
			}
			next(u, w, r)
		}
	}
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
