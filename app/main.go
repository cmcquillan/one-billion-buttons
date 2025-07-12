package main

import (
	"log"
	"net/http"
	"os"

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

	router.UseH2C = true
	router.ForwardedByClientIP = true
	router.SetTrustedProxies(nil)

	buttonApi := ButtonApi{Database: db}
	cursorApi := CursorApi{}

	router.POST("/api/:x/:y", buttonApi.HandlePostButton)

	router.GET("/api/:x/:y", buttonApi.HandleGetButtonPage)

	router.GET("/api/:x/:y/:hash", buttonApi.HandleGetButtonPage)

	router.GET("/cursor/:hex/cursor.png", cursorApi.GetCursor)

	router.GET("/healthcheck/live", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	router.GET("/", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			if err := pusher.Push("/app.js", nil); err != nil {
				log.Printf("Failed on js server push: %v", err)
			}

			if err := pusher.Push("/style.css", nil); err != nil {
				log.Printf("Failed on css server push: %v", err)
			}
		}

		c.Status(http.StatusOK)
		c.File("./static/index.html")
	})
	//router.StaticFile("/", "./static/index.html")
	router.StaticFile("/app.js", "./static/app.js")
	router.StaticFile("/style.css", "./static/style.css")

	router.Run(":8080")
}
