package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var initOk bool = false

var db rsvpDatabase

var templates map[string]*template.Template

var templateFuncs = template.FuncMap{
	"fmtDate": fmtDate,
}

func fmtDate(date time.Time) string {
	return fmt.Sprintf("%d.%d.%d %d:%d:%d",
		date.Day(),
		date.Month(),
		date.Year(),
		date.Hour(),
		date.Minute(),
		date.Second())
}

func init() {
	fmt.Println("Initialize database")
	db = rsvpDatabase{ConnectionString: "./rsvp.db"}
	err := db.create()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = db.initialize()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Initialize templates")
	templates = make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("templates/home.gohtml", "templates/_base.gohtml"))
	templates["form"] = template.Must(template.ParseFiles("templates/form.gohtml", "templates/_base.gohtml"))
	templates["thanks"] = template.Must(template.ParseFiles("templates/thanks.gohtml", "templates/_base.gohtml"))
	templates["list"] = template.Must(template.New("list").Funcs(templateFuncs).ParseFiles("templates/list.gohtml", "templates/_base.gohtml"))
	templates["404"] = template.Must(template.ParseFiles("templates/404.gohtml", "templates/_base.gohtml"))
	templates["500"] = template.Must(template.ParseFiles("templates/500.gohtml", "templates/_base.gohtml"))

	initOk = true
}

func main() {
	if !initOk {
		log.Println("Initialization failed")
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", staticFilesHandler())
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/rsvp", handleForm)
	mux.HandleFunc("/list", basicAuth(handleList))

	server := &http.Server{
		Addr:    ":45612",
		Handler: mux,
	}

	log.Printf("Start listen on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
