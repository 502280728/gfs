// io
package fs

import (
	"fmt"
	"gfs/common"
	"gfs/gfsdk/conf"
)

//仿照HDFS的DFSOutputStream，先将数据分块，然后启动另一个goroutine将数据库存入datanode
//当存入失败的时候，将该数据缓存进入重试队列
type GFSWriter struct {
	buffer        chan *BlockWrapper //数据库的缓存队列，默认容量是10，数据先分成blocksize大小的数据块，缓存进入该channel
	FileName      string             //文件名称，需要这个字段的原因是当一个block准备好的时候，需要向master申请datanode，这时，master需要该值
	totalSize     int64              //总的写入数据量
	currentBuffer *BlockWrapper      //当前需要写入的数据库，大小为blocksize；当写满的时候，将该值推入buffer;然后重新分配blocksize大小的内存
	retryBuffer   chan *BlockWrapper //重试队列，默认容量是10;当重试仍然有失败的时候，整个文件写入失败
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
			data:       make([]byte, 0, conf.GFSContext.GetConf().BlockSize),
			remainSize: conf.GFSContext.GetConf().BlockSize,
		},
	}
	gfs.beginStreamer()
	return gfs
}

func (gfs *GFSWriter) beginStreamer() {
	go func() {
		for v := range gfs.buffer {
			//通知master获取datanode
			//往datanode写数据
			fmt.Println(string(v.data))
		}
	}()
}

func (gfs *GFSWriter) Write(p []byte) (int, error) {
	blockLength := conf.GFSContext.GetConf().BlockSize
	if gfs.currentBuffer.remainSize >= len(p) {
		gfs.currentBuffer.data = append(gfs.currentBuffer.data, p...)
		gfs.totalSize = gfs.totalSize + int64(len(p))
		gfs.currentBuffer.block = getBlockIndex(gfs.totalSize)
		gfs.currentBuffer.fileName = gfs.FileName
		if len(gfs.currentBuffer.data) == blockLength {
			gfs.buffer <- gfs.currentBuffer
		}
	} else {
		loop := (len(p) - gfs.currentBuffer.remainSize) / blockLength
		locate := gfs.currentBuffer.remainSize
		gfs.currentBuffer.data = append(gfs.currentBuffer.data, p[0:gfs.currentBuffer.remainSize]...)
		gfs.totalSize = gfs.totalSize + int64(gfs.currentBuffer.remainSize)
		gfs.currentBuffer.block = getBlockIndex(gfs.totalSize)
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
			gfs.currentBuffer.block = getBlockIndex(gfs.totalSize,blockLength)
			if len(gfs.currentBuffer.data) == blockLength {
				gfs.buffer <- gfs.currentBuffer
			}
		}
	}
	return len(p), nil
}

func getBlockIndex(total int64,blockLength int) int {
	blockLength:=
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
	return nil
}

func (gfs *GFSWriter) Flush() error {
	return nil
}
