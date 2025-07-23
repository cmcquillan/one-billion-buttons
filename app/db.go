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
	Order    int64  `json:"order"`
}

type ObbDb interface {
	GetPageButtonState(x int64, y int64) ([]byte, error)
	SetButtonState(x int64, y int64, index int64, rgb []byte) error
	GetButtonStats() ([]ButtonStat, error)
	LogButtonEvent(x uint64, y uint64, id int64, eventType ButtonEventType) error
	AdjustStat(statKey string, delta int64) error
	RefreshStats() error
}

type ObbDbSql struct {
	connStr string
}

func openConnAndExec(db *ObbDbSql, exec func(dbc *sql.DB) error) error {
	dbc, err := sql.Open("postgres", db.connStr)

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

func (db *ObbDbSql) RefreshStats() error {
	err := openConnAndExec(db, func(dbc *sql.DB) error {
		_, err := dbc.Exec("call update_button_stats()")
		return err
	})

	return err
}

func (db *ObbDbSql) AdjustStat(statKey string, delta int64) error {
	err := openConnAndExec(db, func(dbc *sql.DB) error {
		_, err := dbc.Exec("update button_stat set val = val + $1 where stat_key = $2", delta, statKey)
		return err
	})

	return err
}

func (db *ObbDbSql) LogButtonEvent(x uint64, y uint64, id int64, eventType ButtonEventType) error {
	err := openConnAndExec(db, func(dbc *sql.DB) error {
		_, err := dbc.Exec("insert into button_event (x_coord, y_coord, button_id, event_type) values ($1, $2, $3, $4)", x, y, id, eventType)
		return err
	})

	return err
}

func (db *ObbDbSql) GetButtonStats() ([]ButtonStat, error) {
	rows := 0

	var res *sql.Rows

	err := openConnAndExec(db, func(dbc *sql.DB) error {
		iRes, err := dbc.Query("select stat_key, stat_name, stat_desc, val, scale, \"order\" from button_stat")
		res = iRes
		return err
	})

	result := make([]ButtonStat, 10)

	for res.Next() {
		result[rows] = ButtonStat{}

		res.Scan(&result[rows].StatKey,
			&result[rows].StatName,
			&result[rows].StatDesc,
			&result[rows].Val,
			&result[rows].Scale,
			&result[rows].Order)
		rows++
	}

	return result[0:rows], err
}

func (db *ObbDbSql) GetPageButtonState(x int64, y int64) ([]byte, error) {
	var rows *sql.Rows

	err := openConnAndExec(db, func(dbc *sql.DB) error {
		iRows, err := dbc.Query("select buttons from button where x_coord = $1 and y_coord = $2;", x, y)
		rows = iRows
		return err
	})

	if rows.Next() {
		var bytes []byte
		if serr := rows.Scan(&bytes); err != nil {
			log.Printf("could not scan button: %v", err)
			return nil, serr
		}

		return bytes, err
	}

	return nil, errors.New("coordinate not found")
}

func (db *ObbDbSql) SetButtonState(x int64, y int64, index int64, rgb []byte) error {
	err := openConnAndExec(db, func(dbc *sql.DB) error {
		stmt, err := dbc.Prepare("call set_button_color ($1, $2, $3, $4)")

		if err != nil {
			log.Fatal(err)
			return err
		}

		defer stmt.Close()

		log.Printf("setting (%d, %d, %d) to %s", x, y, index, rgb)
		_, err2 := stmt.Exec(x, y, index, rgb)
		return err2
	})

	return err
}
