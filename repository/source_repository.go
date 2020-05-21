package repository

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/miun173/rss-reader/model"
	"github.com/patrickmn/go-cache"
)

const (
	_defCacherExpiry = time.Minute * 5
)

// SourceRepository :nodoc:
type SourceRepository struct {
	db     *gorm.DB
	cacher *cache.Cache
}

// NewSourceRepository :nodoc:
func NewSourceRepository(db *gorm.DB) *SourceRepository {
	return &SourceRepository{
		db:     db,
		cacher: cache.New(_defCacherExpiry, 1),
	}
}

// FindByID :nodoc:
func (s *SourceRepository) FindByID(id int64) (*model.Source, error) {
	cacheKey := model.NewSourceCacheKeyByID(id)
	src, err := s.findFromCache(cacheKey)
	if err != nil {
		log.Println("error : ", err)
		return nil, err
	}

	if src != nil {
		return src, nil
	}

	source := model.Source{}
	err = s.db.Where("id = ?", id).
		First(&source).
		Error

	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		log.Println("error : ", err)
		return nil, err
	}

	// cache it
	s.cacher.Set(cacheKey, source, _defCacherExpiry)

	return &source, nil
}

// FindAll :nodoc:
func (s *SourceRepository) FindAll(size, page int) (sources []model.Source, err error) {
	if size <= 0 || size > _maxQuerySize {
		size = _maxQuerySize
	}

	if page < 0 {
		page = 0
	}

	err = s.db.
		Limit(size).
		Offset(page).
		Find(&sources).
		Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Println("error: ", err)
		return sources, err
	}

	return sources, nil
}

func (s *SourceRepository) findFromCache(cacheKey string) (*model.Source, error) {
	item, ok := s.cacher.Get(cacheKey)
	if !ok {
		log.Println("not found in cache")
		return nil, nil
	}

	src, ok := item.(model.Source)
	if !ok {
		return nil, errors.New("failed to cast")
	}

	log.Println("found in cache")

	return &src, nil
}
