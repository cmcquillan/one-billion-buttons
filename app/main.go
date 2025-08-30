package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cmcquillan/one-billion-buttons/dblib"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetOutput(os.Stdout)

	appPath, err := os.Executable()

	if err != nil {
		log.Fatal("could not find executable path")
		os.Exit(1)
	}

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Printf("running app from directory: %s", appPath)

	db := &ObbDbSql{
		connStr: cfg.PgConnectionString,
	}

	locker := &dblib.LockSql{
		ConnStr: cfg.PgConnectionString,
	}

	mmDb := &MinimapDbSql{
		connStr: cfg.PgConnectionString,
	}

	ctx, cancel := context.WithCancel(context.Background())

	buttonEventChannel := make(chan BackgroundButtonEvent, cfg.ButtonEventChannelSize)
	go BackgroundEventHandler(db, buttonEventChannel, cfg)
	go BackgroundComputeStatistics(db, ctx, cfg)

	if cfg.RunMinimapInMain {
		log.Print("Starting minimap generation in main instance")
		go BackgroundWorkerMinimap(locker, mmDb, ctx, cfg)
	}

	router := gin.Default()

	router.UseH2C = true
	router.ForwardedByClientIP = true
	router.SetTrustedProxies(nil)

	buttonApi := ButtonApi{Database: db, EventChannel: buttonEventChannel}
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

	router.GET("/minimap.png", func(c *gin.Context) {
		if _, err := os.Stat("./static/minimap.png"); err == nil {
			c.File("./static/minimap.png")
			return
		}

		c.Status(http.StatusNotFound)
	})

	statsApi := StatsApi{Database: db}
	router.GET("/api/stats", statsApi.HandleGetButtonStats)

	router.StaticFile("/app.js", "./static/app.js")
	router.StaticFile("/style.css", "./static/style.css")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Print("Starting one billion buttons http server")
		if err := server.ListenAndServe(); err != nil {
			log.Printf("%v", err)
		}
	}()

	sigChan := make(chan os.Signal, cfg.SignalChannelSize)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	<-sigChan

	cancel()
	log.Print("Shutting down one billion buttons http server")
	server.Shutdown(context.TODO())

	close(buttonEventChannel)
}
