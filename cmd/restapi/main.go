package main

import (
	"log"

	"github.com/miun173/rss-reader/config"
	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/repository"
	"github.com/miun173/rss-reader/restapi"
	"github.com/valyala/fasthttp"
)

func main() {
	cfg := config.GetConfig()

	dbConn := db.NewSQLite3()
	db.Migrate(dbConn)
	defer dbConn.Close()

	sourceRepo := repository.NewSourceRepository(dbConn)
	rssItemRepo := repository.NewRSSItemRepository(dbConn)

	server := restapi.NewServer(sourceRepo, rssItemRepo)
	log.Printf("server start @ :%s\n", cfg.Port())
	log.Fatal(fasthttp.
		ListenAndServe(":"+cfg.Port(), server.Router().Handler),
	)
}
