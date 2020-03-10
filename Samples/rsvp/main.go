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
var tmpl *template.Template

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

	templateFuncs := template.FuncMap{
		"fmtDate": fmtDate,
	}

	tmpl = template.Must(template.New("").Funcs(templateFuncs).ParseGlob("./templates/*"))

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
		Addr:    ":5000",
		Handler: mux,
	}

	log.Printf("Start listen on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
