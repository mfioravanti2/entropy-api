package sqlite3

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	ENGINE = "sqlite3"
)

func Open( connectString string ) ( *gorm.DB, error ) {
	db, err := gorm.Open( ENGINE, connectString )

	return db, err
}
