package repository

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/miun173/rss-reader/model"
)

// RSSItemRepository :nodoc:
type RSSItemRepository struct {
	db *gorm.DB
}

// NewRSSItemRepository :nodoc:
func NewRSSItemRepository(db *gorm.DB) *RSSItemRepository {
	return &RSSItemRepository{
		db: db,
	}
}

// SaveMany :nodoc:
func (r *RSSItemRepository) SaveMany(items []model.RSSItem) {
	for _, item := range items {
		oldRSS, err := r.FindByLink(item.Link)
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

		err = r.db.Create(&item).Error
		if err != nil {
			log.Fatal(err)
		}
	}
}

// FindByLink :nodoc:
func (r *RSSItemRepository) FindByLink(link string) (rssItem *model.RSSItem, err error) {
	item := model.RSSItem{}
	err = r.db.First(&item, "link = ?", link).Error
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

// FetchFromSource :nodoc:
func (r *RSSItemRepository) FetchFromSource() []model.RSSItem {
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
