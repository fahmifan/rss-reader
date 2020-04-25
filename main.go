package main

import (
	"log"

	"github.com/jasonlvhit/gocron"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/repository"
)

func main() {
	dbConn := db.NewSQLite3()
	db.Migrate(dbConn)
	defer dbConn.Close()

	rssItemRepo := repository.NewRSSItemRepository(dbConn)

	log.Println("start worker ...")
	count := 1
	gocron.Every(5).Seconds().Do(func() {
		log.Println("fetch >>> ", count)
		items := rssItemRepo.FetchFromSource()
		rssItemRepo.SaveMany(items)
		log.Println("finished")
		count++
	})
	<-gocron.Start()
}
