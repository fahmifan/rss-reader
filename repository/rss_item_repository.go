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
func (r *RSSItemRepository) SaveMany(sourceID int64, items []model.RSSItem) error {
	for _, item := range items {
		oldRSS, err := r.FindByLink(item.Link)
		if err != nil {
			log.Println("error: ", err)
			return err
		}

		// skip if exists
		if oldRSS != nil {
			continue
		}

		err = r.Create(sourceID, &item)
		if err != nil {
			log.Println("error: ", err)
			return err
		}
	}

	return nil
}

// Create :nodoc:
func (r *RSSItemRepository) Create(sourceID int64, item *model.RSSItem) (err error) {
	now := time.Now()
	item.ID = time.Now().UnixNano()
	item.SourceID = sourceID
	item.CreatedAt = now
	item.UpdatedAt = now

	err = r.db.Create(&item).Error
	if err != nil {
		log.Println("error: ", err)
		return err
	}

	return nil
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
func (r *RSSItemRepository) FetchFromSource(source string) ([]model.RSSItem, error) {
	// fetch from API
	res, err := http.Get(source)
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}

	bt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	rss := model.RSS{}
	err = xml.Unmarshal(bt, &rss)
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}

	return rss.Channel.Items, nil
}

// FindBySourceID :nodoc:
func (r *RSSItemRepository) FindBySourceID(sourceID int64, size, page int) ([]model.RSSItem, error) {
	if size <= 0 || size > _maxQuerySize {
		size = _maxQuerySize
	}

	if page < 0 {
		page = 0
	}

	var items []model.RSSItem
	err := r.db.Where("source_id = ?", sourceID).
		Limit(size).
		Offset(page).
		Find(&items).
		Error

	if err != nil {
		log.Println("error : ", err.Error())
		return nil, err
	}

	return items, nil
}
