package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type ObbDb interface {
	GetPageButtonState(x int64, y int64) ([]byte, error)
	SetButtonState(x int64, y int64, index int64, rgb []byte) error
}

type ObbDbSql struct {
	connStr string
}

func (db *ObbDbSql) GetPageButtonState(x int64, y int64) ([]byte, error) {
	dbc, err := sql.Open("postgres", db.connStr)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer dbc.Close()

	rows, err := dbc.Query("select buttons from button where x_coord = $1 and y_coord = $2;", x, y)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if rows.Next() {
		var bytes []byte
		if err := rows.Scan(&bytes); err != nil {
			log.Fatal(err)
			return nil, err
		}

		return bytes, nil
	}

	return nil, errors.New("coordinate not found")
}

func (db *ObbDbSql) SetButtonState(x int64, y int64, index int64, rgb []byte) error {
	dbc, err := sql.Open("postgres", db.connStr)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer dbc.Close()

	stmt, err := dbc.Prepare("call set_button_color ($1, $2, $3, $4)")

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer stmt.Close()

	log.Printf("setting (%d, %d, %d) to %s", x, y, index, rgb)
	_, err2 := stmt.Exec(x, y, index, rgb)

	if err2 != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
