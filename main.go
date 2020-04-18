package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miun173/rss-reader/model"
)

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

var schema = `
	CREATE TABLE IF NOT EXISTS sources (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL
	);

	CREATE TABLE IF NOT EXISTS rss_items (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		link TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL
	);

	CREATE UNIQUE INDEX IF NOT EXISTS rss_items_link_idx ON rss_items(link);
`

func main() {
	dbConn := connectDB()
	defer dbConn.Close()

	items := fetchRSSItems()

	for _, item := range items {
		now := time.Now()
		item.ID = time.Now().UnixNano()
		item.CreatedAt = now
		item.UpdatedAt = now

		err := dbConn.Create(&item).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("finished")
}

func connectDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "rss-reader.db")
	if err != nil {
		panic("failed to connect database")
	}

	err = db.Exec(schema).Error
	if err != nil {
		log.Fatal(err)
	}

	return db
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
