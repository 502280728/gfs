// config.go
package data

import (
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gfs/gfsnode/data")

var LocalStore *FileStore
