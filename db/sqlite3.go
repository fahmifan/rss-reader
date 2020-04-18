package db

import (
	"log"

	"github.com/jinzhu/gorm"
	// driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewSQLite3 :nodoc:
func NewSQLite3() *gorm.DB {
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
