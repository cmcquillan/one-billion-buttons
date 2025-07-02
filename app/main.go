package main

import (
	b64 "encoding/base64"
	"image"
	"image/png"
	"log"
	"net/http"
	url "net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetOutput(os.Stdout)

	appPath, err := os.Executable()

	if err != nil {
		log.Fatal("could not find executable path")
		os.Exit(1)
	}

	connStr := os.Getenv("PG_CONNECTION_STRING")

	if len(connStr) == 0 {
		log.Fatal("required environment: PG_CONNECTION_STRING")
		os.Exit(1)
	}

	log.Printf("running app from directory: %s", appPath)

	db := &ObbDbSql{
		connStr: connStr,
	}

	router := gin.Default()

	router.POST("/api/:x/:y", func(c *gin.Context) {
		xCoord, errX := strconv.ParseInt(c.Param("x"), 10, 64)
		yCoord, errY := strconv.ParseInt(c.Param("y"), 10, 64)

		if errX != nil || errY != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Could not extract coordinate",
			})

			return
		}

		dto := ButtonStateDto{}
		err := c.BindJSON(&dto)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ix, err := ButtonLocationToIndex(xCoord, yCoord, dto.ID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		rgb, err := HexToBytes(dto.Hex)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		err = db.SetButtonState(xCoord, yCoord, ix, rgb)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Could not complete request",
			})

			return
		}

		bDto, _ := retrieveAndMapGridCoordinate(db, xCoord, yCoord, c.Request)

		button := bDto.Buttons[ix]

		res := http.StatusConflict

		if HexCodesAreEquivalent(button.Hex, dto.Hex) {
			res = http.StatusOK
		}

		c.JSON(res, bDto)
	})

	router.GET("/api/:x/:y", func(c *gin.Context) {
		xCoord, errX := strconv.ParseInt(c.Param("x"), 10, 64)
		yCoord, errY := strconv.ParseInt(c.Param("y"), 10, 64)

		if errX != nil || errY != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Could not extract coordinate",
			})

			return
		}

		if xCoord <= 0 || yCoord <= 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No coordinates at or below 0",
			})

			return
		}

		dto, err := retrieveAndMapGridCoordinate(db, xCoord, yCoord, c.Request)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "We are not available",
			})
		}

		c.JSON(http.StatusOK, dto)
	})

	router.GET("/cursor/:hex/cursor.png", func(c *gin.Context) {
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
	})

	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/app.js", "./static/app.js")
	router.StaticFile("/style.css", "./static/style.css")

	router.Run(":8080")
}

func retrieveAndMapGridCoordinate(db *ObbDbSql, xCoord int64, yCoord int64, r *http.Request) (*GridPageDto, error) {
	state, err := db.GetPageButtonState(xCoord, yCoord)

	if err != nil {
		return nil, err
	}

	page := CreateGridPage(xCoord, yCoord, state)

	data := make([]ButtonStateDto, len(page.Buttons))

	for i, dx := 0, 0; i < len(data); i, dx = i+1, dx+3 {
		var b string = ""
		if !page.Buttons[i].IsEmpty() {
			b = page.Buttons[i].ToHex()
		}
		data[i] = ButtonStateDto{
			Hex: b,
			ID:  page.Buttons[i].id,
		}
	}

	nextHash := b64.StdEncoding.EncodeToString(state)

	nextUri := url.URL{
		Scheme:   r.URL.Scheme,
		Host:     r.URL.Host,
		Path:     r.URL.Path,
		RawQuery: "v=" + nextHash,
	}

	dto := &GridPageDto{
		X:       xCoord,
		Y:       yCoord,
		Buttons: data,
		Next:    nextUri.String(),
	}

	return dto, nil
}
