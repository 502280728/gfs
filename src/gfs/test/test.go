// test
package main

import (
	"gfs/gfsmaster/common"
	"gfs/gfsmaster/fs"
)

func main() {
	mcommon.LoadConf("d:/temp/conf/gfs.properties")
	//	fs.StoreImage()
	fs.RestoreFromLocal()
}
