package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/repository"
	"github.com/miun173/rss-reader/worker"
)

func main() {
	dbConn := db.NewSQLite3()
	db.Migrate(dbConn)
	defer dbConn.Close()
	sech := gocron.NewScheduler(time.Local)

	rssItemRepo := repository.NewRSSItemRepository(dbConn)
	sourceRepo := repository.NewSourceRepository(dbConn)
	wrk := worker.NewWorker(sourceRepo, rssItemRepo)

	log.Println("start worker ...")
	sech.Every(1).Minute().Do(func() {
		wrk.FetchRSS()
		log.Println("finished")
	})
	sech.StartBlocking()
}
