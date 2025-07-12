package main

import (
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type CursorApi struct {
}

func (api *CursorApi) GetCursor(c *gin.Context) {
	hex := c.Param("hex")
	reader, err := os.Open("static/cursor.png")

	if err != nil {
		log.Fatal(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	defer reader.Close()

	img, _, err := image.Decode(reader)

	if err != nil {
		log.Fatal(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	cursor, err := ColorizeImage(img, hex)

	if err != nil {
		log.Fatal(err)
		c.Status(http.StatusInternalServerError)
	}

	errWr := png.Encode(c.Writer, cursor)

	if errWr != nil {
		log.Fatal(errWr)
		c.Status(http.StatusInternalServerError)
	}

	c.Status(http.StatusOK)
}
