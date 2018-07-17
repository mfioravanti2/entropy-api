package sqlite3

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	ENGINE = "sqlite3"
)

// Return an SQLite3 database object
// Set the active status so the service attempts to connect and process objects
func Open( connectString string ) ( *gorm.DB, bool, error ) {
	db, err := gorm.Open( ENGINE, connectString )

	return db, true, err
}
