package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/miun173/rss-reader/model"
)

// SourceRepository :nodoc:
type SourceRepository struct {
	db *gorm.DB
}

// NewSourceRepository :nodoc:
func NewSourceRepository(db *gorm.DB) *SourceRepository {
	return &SourceRepository{db: db}
}

// FindByID :nodoc:
func (s *SourceRepository) FindByID(id int64) (*model.Source, error) {
	source := model.Source{}
	err := s.db.Where("id = ?", id).
		First(&source).
		Error

	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		log.Println("error : ", err)
		return nil, err
	}

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
