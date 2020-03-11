package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

var t *template.Template

func main() {

	t = template.Must(template.New("").ParseGlob("./templates/*"))

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
	err := t.ExecuteTemplate(w, "home", u)
	if err != nil {
		log.Fatal(err)
	}
}

func signinHandler(u user, w http.ResponseWriter, r *http.Request) {
	if u.Authenticated {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	switch r.Method {
	case "GET":
		err := t.ExecuteTemplate(w, "signin", nil)
		if err != nil {
			log.Fatal(err)
		}
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		remember := r.FormValue("rememberme") != ""

		if username != "Admin" || password != "Password123" {
			http.Error(w, "Wrong username or password", http.StatusUnauthorized)
			return
		}

		signinCookie := &http.Cookie{
			Name:     "signin",
			Value:    username,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		}

		if remember {
			signinCookie.MaxAge = 365 * 24 * 60 * 60
		}

		http.SetCookie(w, signinCookie)

		redirectURL := r.URL.Query().Get("returnURL")
		if redirectURL == "" {
			redirectURL = "/"
		}

		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}

func signoutHandler(u user, w http.ResponseWriter, r *http.Request) {
	signoutCookie := &http.Cookie{
		Name:    "signin",
		Value:   "",
		Expires: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Now().Location()),
	}

	http.SetCookie(w, signoutCookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func secureHandler(u user, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Secure site. Welcome %s!", u.UserName)
}
