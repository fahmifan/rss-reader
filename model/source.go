package model

import (
	"fmt"
	"time"
)

// Source :nodoc:
type Source struct {
	ID        int64      `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

// NewSourceCacheKeyByID :nodoc:
func NewSourceCacheKeyByID(id int64) string {
	return fmt.Sprintf("source.id.%d", id)
}
