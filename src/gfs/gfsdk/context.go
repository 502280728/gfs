// context 系统上下文，持有系统配置，文件系统等等
package gfsdk

import (
	"gfs/common"
	"gfs/common/gfs"
)

type GFSContext struct {
	common.Context
}

func GetContext() *GFSContext {
	var gc GFSContext
	gc.RegisterFileSystem("go-filesystem", &GoFileSystem{Name: "go-filesystem"})
	return &gc
}

func (gc *GFSContext) GetDefaultFileSystem() gfs.FileSystem {
	return gc.GetFileSystem("go-filesystem")
}

func (gc *GFSContext) LoadConf(confFile ...string) {
	gc.Context.LoadConf(confFile...)
	gfsc.load()
}
