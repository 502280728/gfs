// io专门为gfs工作的gio实现
package gfsdk

import (
	"sync"
)

//仿照HDFS的DFSOutputStream，先将数据分块，然后启动另一个goroutine将数据块存入datanode
//当存入失败的时候，将该数据缓存进入重试队列
type GFSWriter struct {
	wg              sync.WaitGroup     //用于与streamer同步,即仅当streamer完成工作之后，该writer才能被close，否则会造成主goroutine已经关闭，但是streamer还没有完成工作的情况
	buffer          chan *BlockWrapper //数据库的缓存队列，默认容量是10，数据先分成blocksize大小的数据块，缓存进入该channel
	FileName        string             //文件名称，需要这个字段的原因是当一个block准备好的时候，需要向master申请datanode，这时，master需要该值
	totalSize       int64              //总的写入数据量
	currentBuffer   *BlockWrapper      //当前需要写入的数据库，大小为blocksize；当写满的时候，将该值推入buffer;然后重新分配blocksize大小的内存
	isStreamerStart bool               //表示该流对应的streamer是否启动,为了节约资源，在开始调用Write的地方才启动对应的streamer
}

//代表写入文件时，在client端分割的一个文件块，data大小为blocklength
type BlockWrapper struct {
	data       []byte //一个文件块的数据，最大大小为blocklength
	block      int    //第几个文件块
	fileName   string //文件的名称
	remainSize int    //currentBuffer中剩余的空间大小
}

func NewGFSWriter(fileName string) *GFSWriter {
	gfs := &GFSWriter{
		FileName: fileName,
		buffer:   make(chan *BlockWrapper, 10),
		currentBuffer: &BlockWrapper{
			data:       make([]byte, 0, gfsc.blockSize),
			remainSize: gfsc.blockSize,
		},
	}
	return gfs
}

func (gfs *GFSWriter) beginStreamer(id int) {
	gs := &gfsStreamer{wg: &gfs.wg, buffer: gfs.buffer, id: id}
	gs.start()
}

func (gfs *GFSWriter) Write(p []byte) (int, error) {
	if !gfs.isStreamerStart {
		for i := 0; i < gfsc.streamerSize; i++ {
			gfs.beginStreamer(i)
		}
	}
	blockLength := gfsc.blockSize
	if gfs.currentBuffer.remainSize >= len(p) {
		gfs.currentBuffer.data = append(gfs.currentBuffer.data, p...)
		gfs.totalSize = gfs.totalSize + int64(len(p))
		gfs.currentBuffer.block = getBlockIndex(gfs.totalSize, blockLength)
		gfs.currentBuffer.fileName = gfs.FileName
		if len(gfs.currentBuffer.data) == blockLength {
			gfs.buffer <- gfs.currentBuffer
		}
	} else {
		loop := (len(p) - gfs.currentBuffer.remainSize) / blockLength
		locate := gfs.currentBuffer.remainSize
		gfs.currentBuffer.data = append(gfs.currentBuffer.data, p[0:gfs.currentBuffer.remainSize]...)
		gfs.totalSize = gfs.totalSize + int64(gfs.currentBuffer.remainSize)
		gfs.currentBuffer.block = getBlockIndex(gfs.totalSize, blockLength)
		gfs.currentBuffer.fileName = gfs.FileName
		gfs.buffer <- gfs.currentBuffer
		for i := 0; i <= loop; i++ {
			begin := locate + i*blockLength
			end := begin + blockLength
			if end > len(p) {
				end = len(p)
			}
			gfs.currentBuffer = &BlockWrapper{
				data:       p[begin:end],
				fileName:   gfs.FileName,
				remainSize: blockLength - (end - begin),
			}
			gfs.totalSize = gfs.totalSize + int64(end-begin)
			gfs.currentBuffer.block = getBlockIndex(gfs.totalSize, blockLength)
			if len(gfs.currentBuffer.data) == blockLength {
				gfs.buffer <- gfs.currentBuffer
			}
		}
	}
	return len(p), nil
}

func getBlockIndex(total int64, blockLength int) int {
	if int(total) == 0 {
		return 0
	} else if total%int64(blockLength) == 0 {
		return int(total/int64(blockLength)) - 1
	} else {
		return int(total / int64(blockLength))
	}
}

func (gfs *GFSWriter) Close() error {
	gfs.buffer <- gfs.currentBuffer
	close(gfs.buffer)
	gfs.wg.Wait()
	return nil
}

func (gfs *GFSWriter) Flush() error {
	return nil
}
