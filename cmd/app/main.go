package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jasonlvhit/gocron"
	"github.com/miun173/rss-reader/config"
	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/repository"
	"github.com/miun173/rss-reader/restapi"
	"github.com/miun173/rss-reader/worker"
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
	wrk := worker.NewWorker(sourceRepo, rssItemRepo)

	go func() {
		log.Printf("server start @ :%s\n", cfg.Port())
		log.Fatal(fasthttp.
			ListenAndServe(":"+cfg.Port(), server.Router().Handler),
		)
	}()

	go func() {
		log.Printf("start worker ... every %d minute\n", cfg.CronInterval())
		gocron.Every(cfg.CronInterval()).Minute().Do(func() {
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
