package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron"
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
	cron := gocron.NewScheduler(time.Local)

	srv := &fasthttp.Server{
		Handler: server.Router().Handler,
	}
	go func() {
		log.Println("server start @ :8080")
		if err := srv.ListenAndServe(":8080"); err != nil {
			log.Println(err)
		}
	}()

	cron.Every(30).Second().Do(func() {
		wrk.FetchRSS()
	})

	log.Println("start worker ...")
	cron.StartAsync()

	// block main go routine
	osCh := make(chan os.Signal, 1)
	stopCh := make(chan bool)
	signal.Notify(osCh, os.Interrupt)
	go func() {
		<-osCh
		log.Println("exiting process")
		stopCh <- true
	}()

	<-stopCh
	log.Println("stopping cron")
	cron.Stop()
	log.Println("cron stopped")

	log.Println("stopping server")
	if err := srv.Shutdown(); err != nil {
		log.Println(err)
	}
	log.Println("server stopped")
}
