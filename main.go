package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miun173/rss-reader/db"
	"github.com/miun173/rss-reader/model"
)

func main() {
	dbConn := db.NewSQLite3()
	db.Migrate(dbConn)
	defer dbConn.Close()

	items := fetchRSSItems()

	for _, item := range items {
		oldRSS, err := findRSSItemByLink(dbConn, item.Link)
		if err != nil {
			log.Fatal(err)
		}

		// skip if exists
		if oldRSS != nil {
			continue
		}

		now := time.Now()
		item.ID = time.Now().UnixNano()
		item.CreatedAt = now
		item.UpdatedAt = now

		err = dbConn.Create(&item).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("finished")
}

func findRSSItemByLink(db *gorm.DB, link string) (rssItem *model.RSSItem, err error) {
	item := model.RSSItem{}
	err = db.First(&item, "link = ?", link).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rssItem = &item

	return
}

func fetchRSSItems() []model.RSSItem {
	// fetch from API
	res, err := http.Get("https://hnrss.org/newest")
	if err != nil {
		log.Fatal(err)
	}

	bt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	rss := model.RSS{}
	err = xml.Unmarshal(bt, &rss)
	if err != nil {
		log.Fatal(err)
	}

	return rss.Channel.Items
}
