package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jasonlvhit/gocron"
	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/repository"
	"github.com/miun173/rss-reader/restapi"
	"github.com/miun173/rss-reader/worker"
	"github.com/valyala/fasthttp"
)

func main() {
	dbConn := db.NewSQLite3()
	db.Migrate(dbConn)
	defer dbConn.Close()

	sourceRepo := repository.NewSourceRepository(dbConn)
	rssItemRepo := repository.NewRSSItemRepository(dbConn)

	server := restapi.NewServer(sourceRepo, rssItemRepo)
	wrk := worker.NewWorker(sourceRepo, rssItemRepo)

	go func() {
		log.Println("server start @ :8080")
		log.Fatal(fasthttp.
			ListenAndServe(":8080", server.Router().Handler),
		)
	}()

	go func() {
		log.Println("start worker ...")
		gocron.Every(1).Minute().Do(func() {
			wrk.FetchRSS()
			log.Println("finished")
		})
		<-gocron.Start()
	}()

	// block main go routine
	osCh := make(chan os.Signal)
	stopCh := make(chan bool)
	signal.Notify(osCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-osCh
		log.Println("exiting process")
		stopCh <- true
		os.Exit(0)
	}()

	<-stopCh
}
