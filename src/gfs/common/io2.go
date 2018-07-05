// io2
package common

import (
	"fmt"
	"io"
)

type GFSReadWriteCloser struct {
	TargetNode []*FileLocation
	TargetFile string
	location   int64
	MaxSize    int64
}

func (gfs *GFSReadWriteCloser) Write(p []byte) (int, error) {
	logger.Info("start to write")
	var toIndex = 0     //思路是将切片p分割中大小为MaxWrite的子切片，一次写入MaxWrite大小的数据。toIndex表示第几个子切片
	var end = false     //是否结束循环
	var todo = []byte{} //表示写入targetnode的数据，该切片是p的子切片
	//一次循环写入切片p中最大为MaxWrite的子切片
	for {
		if (toIndex+1)*MaxWrite >= len(p) {
			end = true
			todo = p[toIndex*MaxWrite : len(p)]
		} else {
			todo = p[toIndex*MaxWrite : (toIndex+1)*MaxWrite]
		}
		toIndex++

		targetIndex := gfs.location / BlockSize //表示切片todo里的数据需要写入那个targetnode

		a := BlockSize - gfs.location%BlockSize //计算该block还能接收的字节数
		if int(a) < len(todo) {                 //如果todo的长度大于还能接受的长度，那么要把todo分成两段传输，当然这么做的前提是block的大小不能小于maxwrite
			sendDataTo(gfs.TargetNode[targetIndex], gfs.TargetFile, todo[:a], targetIndex)
			sendDataTo(gfs.TargetNode[targetIndex+1], gfs.TargetFile, todo[a:len(todo)], targetIndex+1)
		} else {
			sendDataTo(gfs.TargetNode[targetIndex], gfs.TargetFile, todo, targetIndex)
		}

		gfs.location = gfs.location + int64(len(todo))

		if end || gfs.location >= gfs.MaxSize { //已经写完切片p或者写完整个文件，就退出
			break
		}
	}
	return len(p), nil
}

func (gfs *GFSReadWriteCloser) Read(p []byte) (n int, err error) {
	if cap(p) < MaxRead {
		return 0, fmt.Errorf("the cap of byte slice is smaller than %d,use a bigger one", MaxRead)
	}
	var pIndex int = 0 //指向切片p中开始写入的位置
	for {
		if pIndex >= cap(p) || gfs.location >= gfs.MaxSize { //只要是把切片p写满或者读到了文件末尾，就退出
			break
		}
		targetIndex := int(gfs.location / BlockSize) //从哪个targetnode读取

		//一次可以读取的数据要从切片p剩余的空间、一个block剩余的数据、以及MaxRead中选择最小的那一个;
		//但是就算如此，最终返回的数据量也可能小于该值。因为可能block中没有那么的数据
		remaina := int(BlockSize - gfs.location%BlockSize)
		remainb := int(cap(p) - pIndex)
		remainc := int(MaxRead)

		var remain int
		if remaina <= remainb {
			remain = remaina
		} else {
			remain = remainb
		}
		if remainc < remain {
			remain = remainc
		}

		if bb, err := getDataFrom(gfs.TargetNode[targetIndex], gfs.TargetFile, gfs.location%BlockSize, int64(remain), targetIndex); err == nil {
			copy(p[pIndex:], bb)
			pIndex = pIndex + len(bb)
			gfs.location = gfs.location + int64(len(bb))

		}

	}
	if pIndex > 0 {
		return pIndex, nil
	} else {
		return 0, io.EOF
	}
}
func (gfs *GFSReadWriteCloser) Close() error {
	return nil
}
