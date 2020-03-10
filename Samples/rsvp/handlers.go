package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

type pageModel struct {
	Title string
}

type notFoundModel struct {
	Title string
	Path  string
}

type formModel struct {
	Title            string
	FirstName        string
	LastName         string
	Email            string
	Attend           bool
	ValidationErrors map[string]string
}

func (m *formModel) Valid() bool {
	m.ValidationErrors = make(map[string]string)

	if strings.TrimSpace(m.FirstName) == "" {
		m.ValidationErrors["FirstName"] = "First name is required"
	}

	if strings.TrimSpace(m.LastName) == "" {
		m.ValidationErrors["LastName"] = "Last name is required"
	}

	if strings.TrimSpace(m.Email) == "" {
		m.ValidationErrors["Email"] = "Email is required"
	}

	return len(m.ValidationErrors) == 0
}

type listModel struct {
	Title       string
	Submissions []submission
}

func staticFilesHandler() http.Handler {
	fs := http.FileServer(http.Dir("static/"))
	return http.StripPrefix("/static/", fs)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		view(w, "home", &pageModel{Title: "Super duper party"})
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	data := &formModel{Title: "Form"}

	switch r.Method {
	case "GET":
		view(w, "form", data)
	case "POST":
		data.FirstName = r.FormValue("fname")
		data.LastName = r.FormValue("lname")
		data.Email = r.FormValue("email")
		data.Attend = r.FormValue("attend") != ""

		if data.Valid() {
			err := db.Save(data.FirstName, data.LastName, data.Email, data.Attend)
			if err != nil {
				log.Fatal(err)
				errorHandler(w, r, http.StatusInternalServerError)
			} else {
				data.Title = "Thank you"
				view(w, "thanks", data)
			}
		} else {
			view(w, "form", data)
		}

	default:
		errorHandler(w, r, http.StatusMethodNotAllowed)
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	submissions, err := db.Submissions()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusInternalServerError)
	} else {
		data := &listModel{"Submissions", submissions}
		view(w, "list", data)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, errorCode int) {
	switch errorCode {
	case http.StatusNotFound:
		w.WriteHeader(http.StatusNotFound)
		data := &notFoundModel{"Not Found", r.URL.Path}
		view(w, "404", data)
	case http.StatusInternalServerError:
		w.WriteHeader(http.StatusInternalServerError)
		data := &pageModel{"Internal Server Error"}
		view(w, "500", data)
	default:
		fmt.Fprintf(w, "Error %d", errorCode)
	}
}

func view(w http.ResponseWriter, name string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
	}
}
