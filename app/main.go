package main

import (
	"log"
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

	buttonApi := ButtonApi{Database: db}
	cursorApi := CursorApi{}

	router.POST("/api/:x/:y", buttonApi.HandlePostButton)

	router.GET("/api/:x/:y", buttonApi.HandleGetButtonPage)

	router.GET("/cursor/:hex/cursor.png", cursorApi.GetCursor)

	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/app.js", "./static/app.js")
	router.StaticFile("/style.css", "./static/style.css")

	router.Run(":8080")
}
