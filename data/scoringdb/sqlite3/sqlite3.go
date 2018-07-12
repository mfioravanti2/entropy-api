package sqlite3

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	ENGINE = "sqlite3"
)

func Open() ( *gorm.DB, error ) {
	db, err := gorm.Open( ENGINE, "./requests.db")

	return db, err
}
