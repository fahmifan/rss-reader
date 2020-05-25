package main

import (
	"log"
	"os"

	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/repository"
	"github.com/miun173/rss-reader/restapi"
	"github.com/valyala/fasthttp"
)

var _port string

func init() {
	_port = "8080"
	if val, ok := os.LookupEnv("HTTP_PORT"); ok {
		_port = val
	}
}

func main() {
	dbConn := db.NewSQLite3()
	db.Migrate(dbConn)
	defer dbConn.Close()

	sourceRepo := repository.NewSourceRepository(dbConn)
	rssItemRepo := repository.NewRSSItemRepository(dbConn)

	server := restapi.NewServer(sourceRepo, rssItemRepo)
	log.Printf("server start @ :%s\n", _port)
	log.Fatal(fasthttp.
		ListenAndServe(":"+_port, server.Router().Handler),
	)
}
