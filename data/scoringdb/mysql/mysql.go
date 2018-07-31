package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	ENGINE = "mysql"
)

// Return an MySQL database object
// Set the active status so the service attempts to connect and process objects
func Open( connectString string ) ( *gorm.DB, bool, error ) {
	db, err := gorm.Open( ENGINE, connectString )

	return db, true, err
}
