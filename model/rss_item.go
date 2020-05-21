package model

import (
	"encoding/xml"
	"time"
)

// RSS :nodoc:
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel :nodoc:
type Channel struct {
	XMLName xml.Name  `xml:"channel"`
	Items   []RSSItem `xml:"item"`
}

// RSSItem :nodoc:
type RSSItem struct {
	XMLName xml.Name `xml:"item" json:"-"`

	ID          int64      `gorm:"primary_key" xml:"-" json:"id"`
	SourceID    int64      `xml:"-" json:"sourceID"`
	Title       string     `xml:"title" json:"title"`
	Description string     `xml:"description" json:"description"`
	Link        string     `xml:"link" json:"link"`
	CreatedAt   time.Time  `xml:"-" json:"createdAt"`
	UpdatedAt   time.Time  `xml:"-" json:"updatedAt"`
	DeletedAt   *time.Time `xml:"-" json:"-"`
}
