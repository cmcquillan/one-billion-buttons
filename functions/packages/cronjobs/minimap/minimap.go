package minimap

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/cmcquillan/one-billion-buttons/dblib"
)

func Main(args map[string]interface{}) map[string]interface{} {

	status := "unknown"

	connStr := os.Getenv("PG_CONNETION_STRING")

	locker := &dblib.LockSql{
		ConnStr: connStr,
	}

	db := &ObbDbSql{
		connStr: connStr,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)

	defer cancel()

	CreateMinimap(locker, db, ctx)

	msg := make(map[string]interface{})
	msg["status"] = status
	return msg
}

const MINIMAP_LOCK_TYPE = "minimap_gen"

func CreateMinimap(locker dblib.Lock, db ObbDb, ctx context.Context) {
	log.Print("Background minimap maker started")

	lockVal, err := locker.AcquireLock(MINIMAP_LOCK_TYPE, time.Hour)

	if err == dblib.ErrLockNotAcquired {
		log.Printf("%s lock already acquired, deferring work", MINIMAP_LOCK_TYPE)
		return
	}

	if err != nil {
		log.Printf("error acquiring lock: %v", err)
		return
	}

	log.Printf("lock %s acquired for %s", lockVal.Value, lockVal.Type)

	mmChan := make(chan *MinimapItem, 10000)

	x, y, err := db.GetImageDimensions()

	if err != nil {
		log.Printf("unable to get correct dimensions for map: %v", err)
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
		return
	}

	err = e.Encode(file, minimap)

	if err != nil {
		log.Printf("could not encode png minimap: %v", err)
	}

	defer log.Printf("lock %s released for %s", lockVal.Value, lockVal.Type)
	defer locker.ReleaseLock(lockVal)
	log.Print("Background minimap maker stopped")
}
