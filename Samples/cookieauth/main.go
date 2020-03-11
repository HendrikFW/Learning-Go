package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type user struct {
	Authenticated bool
	UserName      string
}

type handler func(w http.ResponseWriter, r *http.Request)

type authHandler func(u user, w http.ResponseWriter, r *http.Request)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", allowAnonymous(homeHandler))
	mux.HandleFunc("/signin", method(allowAnonymous(signinHandler), "GET", "POST"))
	mux.HandleFunc("/signout", method(secure(signoutHandler), "GET"))
	mux.HandleFunc("/secure", secure(secureHandler))

	server := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func homeHandler(u user, w http.ResponseWriter, r *http.Request) {
	var msg string
	if u.Authenticated {
		msg = fmt.Sprintf("signed in. Hello, %s!", u.UserName)
	} else {
		msg = "not signed in."
	}

	fmt.Fprintf(w, "Homepage, You're %s", msg)
}

func signinHandler(u user, w http.ResponseWriter, r *http.Request) {
	if u.Authenticated {
		fmt.Fprintf(w, "You're already signed in %s!", u.UserName)
		return
	}

	year := 365 * 24 * 60 * 60
	signInCookie := fmt.Sprintf("signin=%s; Max-Age=%d; HttpOnly; Path=/; SameSite=Strict", "Hendrik", year)
	w.Header().Add("Set-Cookie", signInCookie)

	returnURL := r.URL.Query().Get("returnURL")
	if returnURL != "" {
		w.Header().Add("Location", returnURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		fmt.Fprintf(w, "%s Sign in", r.Method)
	}
}

func signoutHandler(u user, w http.ResponseWriter, r *http.Request) {
	signoutCookie := &http.Cookie{
		Name:    "signin",
		Value:   "",
		Expires: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Now().Location()),
	}

	http.SetCookie(w, signoutCookie)
	fmt.Fprintln(w, "You've been signed out!")
}

func secureHandler(u user, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Secure site. Welcome %s!", u.UserName)
}

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
