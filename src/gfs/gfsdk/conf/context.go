// context 系统上下文，持有系统配置，文件系统等等
package conf

import (
	"gfs/common/gfs"
	"gfs/gfsdk/fs"
)

type Context struct {
	conf *GFSClientConf
	fs   gfs.FileSystem
}

var GFSContext *Context

func init() {
	GFSContext = &Context{
		conf: &GFSClientConf{},
		fs:   &fs.GoFileSystem{Name: "go-filesystem"},
	}
}

func (gc *Context) GetFileSystem() gfs.FileSystem {
	return gc.fs
}

func (gc *Context) LoadConf(confFile ...string) {
	gc.conf.load(confFile...)
}

func (gc *Context) AddConf(conf map[string]string) {
	gc.conf.loadMap(conf)
}

func (gc *Context) GetConf() *GFSClientConf {
	return gc.conf
}
