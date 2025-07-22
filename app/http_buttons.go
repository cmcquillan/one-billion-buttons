package main

import (
	b64 "encoding/base64"
	"net/http"
	url "net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ButtonApi struct {
	Database     ObbDb
	EventChannel chan BackgroundButtonEvent
}

func (api *ButtonApi) HandleGetButtonPage(c *gin.Context) {
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

	dto, err := retrieveAndMapGridCoordinate(api.Database, xCoord, yCoord, c.Request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "We are not available",
		})
	}

	c.JSON(http.StatusOK, dto)
}

func (api *ButtonApi) HandlePostButton(c *gin.Context) {
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

	err = api.Database.SetButtonState(xCoord, yCoord, ix, rgb)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not complete request",
		})

		return
	}

	bDto, _ := retrieveAndMapGridCoordinate(api.Database, xCoord, yCoord, c.Request)

	button := bDto.Buttons[ix]

	res := http.StatusConflict

	if HexCodesAreEquivalent(button.Hex, dto.Hex) {
		res = http.StatusOK
	}

	api.EventChannel <- BackgroundButtonEvent{
		X:     uint32(xCoord),
		Y:     uint32(yCoord),
		ID:    dto.ID,
		Event: ButtonEventTypePress,
	}

	c.JSON(res, bDto)
}

func retrieveAndMapGridCoordinate(db ObbDb, xCoord int64, yCoord int64, r *http.Request) (*GridPageDto, error) {
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
