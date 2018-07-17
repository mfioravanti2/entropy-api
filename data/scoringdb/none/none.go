package none

import (
	"github.com/jinzhu/gorm"
)

const (
	ENGINE = "none"
)

// Return an empty data store object, but set the active flag to false so
// the service treats the connection as a pass through object
func Open( connectString string ) ( *gorm.DB, bool, error ) {

	return nil, false, nil
}
