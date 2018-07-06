// io_streamer
package gfsdk

import (
	"fmt"
	"gfs/common/ghttp"
	"sync"
)

type gfsStreamer struct {
	wg          *sync.WaitGroup
	buffer      chan *BlockWrapper
	retryBuffer chan *BlockWrapper
	id          int //唯一标示，一个流中的stremer的id唯一
}

func (gs *gfsStreamer) start() {
	logger.Infof("streamer %d started", gs.id)
	gs.wg.Add(1)
	go gs.transfer()
}

//传输数据
func (gs *gfsStreamer) transfer() {
	var http ghttp.GFSRequest
	for v := range gs.buffer {
		transfer(http, v)
	}
	gs.wg.Done()
}

func transfer(http ghttp.GFSRequest, v *BlockWrapper) {
	//通知master获取datanode
	//往datanode写数据
	err := http.PostObj(gfsc.master, v.data, nil)
	logger.Error(err)
	fmt.Println(string(v.data))
}
