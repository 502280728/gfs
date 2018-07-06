// conf
package conf

import (
	"gfs/common/gconf"
	"strings"
)

type GFSClientConf struct {
	Master    string //必须是完成的包含protocal://host:port的url
	UserName  string
	BlockSize int //文件块的大小
}

func (gcc *GFSClientConf) load(confFile ...string) {
	for _, cf := range confFile {
		if strings.HasSuffix(cf, ".properties") {
			gconf.Conf.LoadProperties(cf)
		} else if strings.HasSuffix(cf, ".yaml") || strings.HasSuffix(cf, ".yml") {
			gconf.Conf.LoadYaml(cf)
		}
	}
	gcc.put()
}

func (gcc *GFSClientConf) loadMap(conf map[string]string) {
	gconf.Conf.AddExtras(conf)
	gcc.put()
}

func (gcc *GFSClientConf) put() {
	if s, err := gconf.Conf.GetURL(gconf.GFS_MASTER); err == nil {
		gcc.Master = s
	}
	gconf.Conf.GetOrDefault(gconf.GFS_BLOCK_SIZE, &gcc.BlockSize, 1000)
}
