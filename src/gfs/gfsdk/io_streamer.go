// io_streamer
package gfsdk

import (
	"fmt"
	"gfs/common/ghttp"
	"gfs/common/gmsg"
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
	msg := &gmsg.MsgToMaster1{Block: v.block, FileName: v.fileName}
	var fl gmsg.MsgToSDK1
	err := http.PostObj(gfsc.master, msg, &fl)

	//往datanode写数据
	logger.Error(err)
	fmt.Println(string(v.data))
}
