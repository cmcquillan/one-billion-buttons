package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/cmcquillan/one-billion-buttons/dblib"
)

type MinimapItem struct {
	X   int64
	Y   int64
	RGB []byte
}

type ObbDb interface {
	BeginMinimapStreaming(stream chan *MinimapItem, ctx context.Context) error
	GetImageDimensions() (x int64, y int64, e error)
}

type ObbDbSql struct {
	connStr string
}

func (db *ObbDbSql) GetConnectionString() string {
	return db.connStr
}

func (db *ObbDbSql) GetImageDimensions() (x int64, y int64, e error) {

	err := dblib.OpenConnAndExec(db, func(dbc *sql.DB) error {
		row := dbc.QueryRow("select max(x_coord), max(y_coord) from button")

		scanErr := row.Scan(&x, &y)
		return scanErr
	})

	return x, y, err
}

func (db *ObbDbSql) BeginMinimapStreaming(stream chan *MinimapItem, ctx context.Context) error {
	err := dblib.OpenConnAndExec(db, func(dbc *sql.DB) error {

		// Yep, query literally everything and stream to the application. Dirty reads are fine
		rows, err := dbc.QueryContext(ctx, "set transaction isolation level read uncommitted; select x_coord, y_coord, map_value from button;")

		if err != nil {
			return err
		}

		defer rows.Close()
		defer close(stream)

		count := 0

		for rows.Next() {
			cErr := ctx.Err()
			if cErr != nil {
				log.Printf("aborting minimap stream due to context error: %v", cErr)
				return cErr
			}

			mmap := MinimapItem{}

			rows.Scan(&mmap.X, &mmap.Y, &mmap.RGB)
			count++
			stream <- &mmap
		}

		log.Printf("scanned %d rows for minimap", count)
		return nil
	})

	return err
}
