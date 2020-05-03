package model

import "time"

// Source :nodoc:
type Source struct {
	ID        int64      `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}
