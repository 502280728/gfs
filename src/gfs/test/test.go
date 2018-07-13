// test
package main

import (
	"gfs/common/logging"
)

func main() {
	logging.Init("d:/temp/log.log", "")
	logger := logging.GetLogger("df")
	logger.Info("tis is ")
}
