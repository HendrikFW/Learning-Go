package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Learning Go: sqlite")

	database, err := sql.Open("sqlite3", ":memory:")
	checkError(err)

	createTableStmt, err := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstName TEXT, lastName text)")
	checkError(err)

	_, err = createTableStmt.Exec()
	checkError(err)

	insertStmt, err := database.Prepare("INSERT INTO people (firstName, lastName) VALUES (?, ?)")
	checkError(err)

	_, err = insertStmt.Exec("John", "Doe")
	checkError(err)

	rows, err := database.Query("SELECT id, firstName, lastName FROM people")
	checkError(err)
	var id int
	var firstName string
	var lastName string
	for rows.Next() {
		rows.Scan(&id, &firstName, &lastName)
		fmt.Printf("%d: %s, %s\n", id, lastName, firstName)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
