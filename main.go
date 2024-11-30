package main

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"

	"example.com/web/chat"
	"example.com/web/router"
)

//go:embed chatui/prod
var server embed.FS

//go:embed database/schema.sql
var schemaSQL string

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Working directory: ", wd)

	// open the sqlite database file in local dev
	db, err := sql.Open("sqlite", wd+"/database/database.db")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running migrations ...")
	_, err = db.Exec(schemaSQL)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// Create the Gin router
	r := gin.Default()
	r.Use(static.Serve("/", static.EmbedFolder(server, "chatui/prod/build")))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})
	r.NoRoute(func(c *gin.Context) {
		log.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		// Send index.html back for react router paths
		c.File("chatui/prod/build/index.html")
	})

	r.POST("/users", func(c *gin.Context) { router.CreateUser(c, db) })
	r.POST("/login", func(c *gin.Context) { router.Login(c, db) })
	// Listing endpoints
	r.GET("/channels", func(c *gin.Context) { router.ListChannels(c, db) })
	r.GET("/messages", func(c *gin.Context) { router.ListMessages(c, db) })
	r.POST("/channels", func(c *gin.Context) { router.CreateChannel(c, db) })
	r.POST("/messages", func(c *gin.Context) { router.CreateMessage(c, db) })

	// WS stuff
	hub := chat.NewHub()
	go hub.Run()
	r.GET("/ws", func(c *gin.Context) {
		chat.ServeWS(hub, c.Writer, c.Request)
	})

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
