package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type rsvpDatabase struct {
	Db               *sql.DB
	ConnectionString string
}

type submission struct {
	ID        int
	Date      time.Time
	FirstName string
	LastName  string
	Email     string
	Attend    bool
}

func (d *rsvpDatabase) create() error {
	database, err := sql.Open("sqlite3", d.ConnectionString)
	if err != nil {
		return err
	}

	err = database.Ping()
	if err != nil {
		return err
	}

	d.Db = database
	return nil
}

func (d *rsvpDatabase) initialize() error {
	createStmt, err := d.Db.Prepare(
		`CREATE TABLE IF NOT EXISTS submissions (
			ID INTEGER PRIMARY KEY, 
			FirstName TEXT,
			LastName TEXT,
			Email TEXT,
			Attend INTEGER,
			Date INTEGER)`)

	if err != nil {
		return err
	}

	_, err = createStmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (d *rsvpDatabase) Submissions() ([]submission, error) {
	submissions := make([]submission, 0)

	rows, err := d.Db.Query(`SELECT ID, FirstName, LastName, Email, Attend, Date FROM submissions`)
	if err != nil {
		return nil, err
	}

	var id int
	var fname string
	var lname string
	var email string
	var attend bool
	var date int64

	for rows.Next() {
		err = rows.Scan(&id, &fname, &lname, &email, &attend, &date)
		if err != nil {
			return nil, err
		}

		submissions = append(submissions, submission{
			ID:        id,
			FirstName: fname,
			LastName:  lname,
			Email:     email,
			Attend:    attend,
			Date:      time.Unix(date, 0),
		})
	}

	return submissions, nil
}

func (d *rsvpDatabase) Save(firstName, lastName, email string, attend bool) error {
	insertStmt, err := d.Db.Prepare(
		`INSERT INTO submissions 
			(FirstName, LastName, Email, Attend, Date) 
			VALUES (?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}

	date := time.Now().Unix()
	_, err = insertStmt.Exec(firstName, lastName, email, attend, date)

	if err != nil {
		return err
	}

	return nil
}
