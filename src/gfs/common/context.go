// context
package common

import (
	"errors"
	"gfs/common/gconf"
	"gfs/common/gfs"
	"strings"
)

type Context struct {
	Conf    gconf.Configuration
	fsCache map[string]gfs.FileSystem
}

var (
	FS_DUPLICATE_ERROR = errors.New("DUPLICATE FS WHEN REGISTER NEW FS")
)

func GetContext() *Context {
	return &Context{Conf: gconf.Conf, fsCache: make(map[string]gfs.FileSystem)}
}

func (gc *Context) GetFileSystem(name string) gfs.FileSystem {
	return gc.fsCache[name]
}

func (gc *Context) RegisterFileSystem(name string, fs gfs.FileSystem) error {
	if gc.fsCache == nil {
		gc.fsCache = make(map[string]gfs.FileSystem)
	}
	if _, found := gc.fsCache[name]; found {
		return FS_DUPLICATE_ERROR
	} else {
		gc.fsCache[name] = fs
		return nil
	}
}

func (gc *Context) LoadConf(confFile ...string) {
	if gc.Conf == nil {
		gc.Conf = gconf.Conf
	}
	for _, cf := range confFile {
		if strings.HasSuffix(cf, ".properties") {
			gc.Conf.LoadProperties(cf)
		} else if strings.HasSuffix(cf, ".yaml") || strings.HasSuffix(cf, ".yml") {
			gc.Conf.LoadYaml(cf)
		}
	}
}

func (gc *Context) AddConf(conf map[string]string) {
	gc.Conf.AddExtras(conf)
}

func (gc *Context) GetConf() gconf.Configuration {
	return gc.Conf
}
