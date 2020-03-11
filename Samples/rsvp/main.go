package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var db rsvpDatabase
var tmpl *template.Template
var adminPw string

func main() {
	app := &cli.App{
		Name:  "Party invites",
		Usage: "runs a local website for party invites",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Value:   5000,
				Usage:   "Port to listen on",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    "connectionString",
				Value:   ":memory:",
				Usage:   "The SQLite connection string to use",
				Aliases: []string{"cs"},
			},
			&cli.StringFlag{
				Name:        "password",
				Value:       "Test123",
				Usage:       "Use username 'Admin' to view submissions",
				Aliases:     []string{"pw"},
				Destination: &adminPw,
			},
		},
		Action: launch,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func launch(c *cli.Context) error {

	db = rsvpDatabase{ConnectionString: c.String("connectionString")}
	fmt.Printf("Using connection '%s'\n", db.ConnectionString)

	err := db.create()
	if err != nil {
		return err
	}

	err = db.initialize()
	if err != nil {
		return err
	}

	templateFuncs := template.FuncMap{
		"fmtDate": fmtDate,
	}

	tmpl, err = template.New("").Funcs(templateFuncs).ParseGlob("./templates/*")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", staticFilesHandler())
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/rsvp", handleForm)
	mux.HandleFunc("/list", basicAuth(handleList))

	addr := fmt.Sprintf(":%d", c.Int("port"))
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Start listen on %s\n", server.Addr)
	return server.ListenAndServe()
}
