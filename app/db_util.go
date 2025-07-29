package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DbString interface {
	GetConnectionString() string
}

func OpenConnAndExec(db DbString, exec func(dbc *sql.DB) error) error {
	dbc, err := sql.Open("postgres", db.GetConnectionString())

	if err != nil {
		log.Printf("could not open database connection: %v", err)
		return err
	}

	defer dbc.Close()

	if err := exec(dbc); err != nil {
		log.Printf("could not execute operation: %v", err)
		return err
	}

	return nil
}
