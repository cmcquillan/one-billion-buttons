package main

import (
	"context"
	"database/sql"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/cmcquillan/one-billion-buttons/dblib"
)

type MinimapItem struct {
	X   int64
	Y   int64
	RGB []byte
}

type MinimapDb interface {
	BeginMinimapStreaming(stream chan *MinimapItem, ctx context.Context) error
	GetImageDimensions() (x int64, y int64, e error)
}

type MinimapDbSql struct {
	connStr string
}

func (db *MinimapDbSql) GetConnectionString() string {
	return db.connStr
}

func (db *MinimapDbSql) GetImageDimensions() (x int64, y int64, e error) {

	err := dblib.OpenConnAndExec(db, func(dbc *sql.DB) error {
		row := dbc.QueryRow("select max(x_coord), max(y_coord) from button")

		scanErr := row.Scan(&x, &y)
		return scanErr
	})

	return x, y, err
}

func (db *MinimapDbSql) BeginMinimapStreaming(stream chan *MinimapItem, ctx context.Context) error {
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

const MINIMAP_LOCK_TYPE = "minimap_gen"

func BackgroundWorkerMinimap(locker dblib.Lock, db MinimapDb, ctx context.Context) {
	log.Print("Background minimap maker started")

	ticker := time.NewTicker(time.Second)

tickerLoop:
	for {
		select {
		case <-ticker.C:
			if CreateMinimap(locker, db, ctx) {
				ticker.Reset(time.Minute * 10)
			}
		case <-ctx.Done():
			break tickerLoop
		}
	}

	log.Print("Background minimap maker stopped")
}

func CreateMinimap(locker dblib.Lock, db MinimapDb, ctx context.Context) bool {

	lockVal, err := locker.AcquireLock(MINIMAP_LOCK_TYPE, time.Minute*10)

	if err == dblib.ErrLockNotAcquired {
		log.Printf("%s lock already acquired, deferring work", MINIMAP_LOCK_TYPE)
		return false
	}

	if err != nil {
		log.Printf("error acquiring lock: %v", err)
		return false
	}

	defer locker.ReleaseLock(lockVal)

	log.Printf("lock %s acquired for %s", lockVal.Value, lockVal.Type)

	mmChan := make(chan *MinimapItem, 10000)

	x, y, err := db.GetImageDimensions()

	if err != nil {
		log.Printf("unable to get correct dimensions for map: %v", err)
		return false
	}

	go db.BeginMinimapStreaming(mmChan, ctx)

	minimap := image.NewRGBA64(image.Rect(0, 0, int(x), int(y)))

	for mmItem := range mmChan {
		alpha := 0

		if mmItem.RGB[0]+mmItem.RGB[1]+mmItem.RGB[2] > 0 {
			alpha = 255
		}

		c := color.RGBA{
			R: uint8(mmItem.RGB[0]),
			G: uint8(mmItem.RGB[1]),
			B: uint8(mmItem.RGB[2]),
			A: uint8(alpha),
		}

		minimap.Set(int(mmItem.X), int(mmItem.Y), c)
	}

	e := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	file, err := os.OpenFile("./static/minimap.png", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("failed to open minimap file: %v", err)
		return false
	}

	err = e.Encode(file, minimap)

	if err != nil {
		log.Printf("could not encode png minimap: %v", err)
		return false
	}

	return true
}
