// fs
package fs

import (
	"fmt"
	"sync"
)

type Replica struct {
	Locate string // 位置
	Valid  bool   //是否是完整的副本,在某些情况下，一个副本可能是损坏的
}

type FileBlock struct {
	ID      int        //数据块id
	Replica []*Replica //分配的副本
}

//master中存放的文件的位置，Location中存放实际的位置，当Master为Block分配存储空间时，先将分配的
//FileLocate存到TempLocate中，key为blockId，当master收到消息说该Block已被存储时，将之从TempLocate
//移至Location
type FileBlocks struct {
	sync.RWMutex
	Location   []*FileBlock
	TempLocate map[int]*FileBlock
}

//分配一个地址
func (fb *FileBlocks) Allocate(f *FileBlock) {
	fb.Lock()
	defer fb.Unlock()
	fb.TempLocate[f.ID] = f
}

//从一个datanode收到一个block的确认消息
func (fb *FileBlocks) ACK(blockId int, l string) {
	fb.Lock()
	defer fb.Unlock()

	fb.Location = append(fb.Location, fb.TempLocate[blockId])
	delete(fb.TempLocate, blockId)
}
