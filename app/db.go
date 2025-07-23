package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type ButtonStat struct {
	StatKey  string `json:"stat_key"`
	StatName string `json:"stat_name"`
	StatDesc string `json:"stat_desc"`
	Val      int64  `json:"val"`
	Scale    int64  `json:"scale"`
}

type ObbDb interface {
	GetPageButtonState(x int64, y int64) ([]byte, error)
	SetButtonState(x int64, y int64, index int64, rgb []byte) error
	GetButtonStats() ([]ButtonStat, error)
	LogButtonEvent(x uint64, y uint64, id int64, eventType ButtonEventType) error
	AdjustStat(statKey string, delta int64) error
}

type ObbDbSql struct {
	connStr string
}

func (db *ObbDbSql) AdjustStat(statKey string, delta int64) error {
	dbc, err := sql.Open("postgres", db.connStr)

	if err != nil {
		log.Printf("could not open database connection: %v", err)
		return err
	}

	defer dbc.Close()

	_, err = dbc.Exec("update button_stat set val = val + $1 where stat_key = $2", delta, statKey)

	if err != nil {
		log.Printf("could not adjust stat %s by %d: %v", statKey, delta, err)
		return err
	}

	return nil
}

func (db *ObbDbSql) LogButtonEvent(x uint64, y uint64, id int64, eventType ButtonEventType) error {
	dbc, err := sql.Open("postgres", db.connStr)

	if err != nil {
		log.Printf("could not open database connection: %v", err)
		return err
	}

	defer dbc.Close()

	_, err = dbc.Exec("insert into button_event (x_coord, y_coord, button_id, event_type) values ($1, $2, $3, $4)", x, y, id, eventType)

	if err != nil {
		log.Printf("could not log button event %v", err)
		return err
	}
	return nil
}

func (db *ObbDbSql) GetButtonStats() ([]ButtonStat, error) {
	rows := 0

	dbc, err := sql.Open("postgres", db.connStr)
	if err != nil {
		log.Printf("could not open database connection: %v", err)
		return nil, err
	}

	defer dbc.Close()

	res, err := dbc.Query("select stat_key, stat_name, stat_desc, val, scale from button_stat")

	if err != nil {
		log.Printf("could not query button stats: %v", err)
		return nil, err
	}

	result := make([]ButtonStat, 10)

	for res.Next() {
		result[rows] = ButtonStat{}

		res.Scan(&result[rows].StatKey, &result[rows].StatName, &result[rows].StatDesc, &result[rows].Val, &result[rows].Scale)
		rows++
	}

	return result[0:rows], nil
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
