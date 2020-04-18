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
	XMLName xml.Name `xml:"item"`

	ID          int64      `gorm:"primary_key" xml:"-"`
	Title       string     `xml:"title"`
	Description string     `xml:"description"`
	Link        string     `xml:"link"`
	CreatedAt   time.Time  `xml:"-"`
	UpdatedAt   time.Time  `xml:"-"`
	DeletedAt   *time.Time `xml:"-"`
}
