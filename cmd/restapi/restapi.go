package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/konrads/go-tagged-articles/pkg/db"
	"github.com/konrads/go-tagged-articles/pkg/handler"
)

func main() {
	restUri := flag.String("rest-uri", "0.0.0.0:8080", "rest uri")
	dbUri := flag.String("db-uri", "NA" /* eg. for postgres: "postgres://gomicros:password@localhost/gomicros?sslmode=disable" */, "db uri")
	flag.Parse()

	log.Printf(`Starting restapi service with params:
	- restUri: %s
	- dbUri:   %s
	`, *restUri, *dbUri)

	var db db.DB = db.NewPostgresDB(dbUri)
	defer db.Close()
	handler := handler.NewHandler(&db)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/articles/:id", handler.GetArticle)
	r.GET("/tags/:tag/:date", handler.GetTagInfos)
	r.POST("/articles", handler.PostArticle)

	r.Run(*restUri)
}
