// gfsdk
package gfsdk

import (
	"gfs/common/gconf"
)

type gfsClientConf struct {
	master           string //必须是完成的包含protocal://host:port的url
	userName         string
	blockSize        int //文件块的大小
	streamerSize     int //单个文件流的写的goroutine的数量，比如一个文件在写的时候可以有多个goroutine同时往不同的datanode写,默认是2
	streamerChipSize int //每个streamer将数据分成streamerChipSize大小的数据库传输给datanode
}

const (
	GFS_CLIENT_WRITER_STREAMER_SIZE      = "gfs.client.writer.streamer.size"
	GFS_CLIENT_WRITER_STREAMER_CHIP_SIZE = "gfs.client.writer.streamer.chip.size"

	GFS_BLOCK_SIZE = "gfs.block.size"
	GFS_MASTER     = "gfs.master"
)

var gfsc gfsClientConf

//从全球的配置中加载gofssdk需要的配置
func (gcc *gfsClientConf) load() {
	if s, err := gconf.Conf.GetURL(GFS_MASTER); err == nil {
		gcc.master = s
	}
	gconf.Conf.GetOrDefault(GFS_BLOCK_SIZE, &gcc.blockSize, 1000)
	gconf.Conf.GetOrDefault(GFS_CLIENT_WRITER_STREAMER_SIZE, &gcc.streamerSize, 2)
	gconf.Conf.GetOrDefault(GFS_CLIENT_WRITER_STREAMER_CHIP_SIZE, &gcc.streamerChipSize, 50)
}
