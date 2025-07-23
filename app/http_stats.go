package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsApi struct {
	Database ObbDb
}

func (api *StatsApi) HandleGetButtonStats(c *gin.Context) {
	stats, err := api.Database.GetButtonStats()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve button stats",
		})
		return
	}

	c.Header("Cache-Control", "max-age=120, public")
	c.JSON(http.StatusOK, stats)
}
