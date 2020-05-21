package main

import (
	"log"

	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/repository"
	"github.com/miun173/rss-reader/restapi"
	"github.com/valyala/fasthttp"
)

func main() {
	dbConn := db.NewSQLite3()
	db.Migrate(dbConn)
	defer dbConn.Close()

	sourceRepo := repository.NewSourceRepository(dbConn)
	rssItemRepo := repository.NewRSSItemRepository(dbConn)

	server := restapi.NewServer(sourceRepo, rssItemRepo)
	log.Println("server start @ :8080")
	log.Fatal(fasthttp.
		ListenAndServe(":8080", server.Router().Handler),
	)
}
