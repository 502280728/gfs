package common

import (
	"bytes"
	"encoding/gob"
)

type MessageInFS struct {
	Success bool
	Data    string
	Msg     string
}

type ACK bool

//由datanode每隔几秒传给master的信息，包含datanode刚刚完成存储的文件信息
type DataNodeIntervalMessage struct {
	Files []string //格式是 文件名全称:blockid
}

type MasterToClientMessage struct {
	Nodes []string
}

//表示某个文件块在整个文件系统中的存储位置
type FileLocation struct {
	Main    string   //主要存储节点
	Replica []string //其他存储节点,当主节点失效后,第一个Replica自动升为主节点，以此类推
}

type FileBlockChip struct {
	Block    int      //该Block的id
	FileName string   //该Block对应的文件名
	Replica  []string //该Block备份位置
	Offset   int64    //Data在该Block中的偏移量
	Limit    int64    //Data的最大长度
	Data     []byte   //该Block中对应的offset，limit的数据
}

func (fbc *FileBlockChip) Encode() *bytes.Buffer {
	var res *bytes.Buffer
	enc := gob.NewEncoder(res)
	enc.Encode(*fbc)
	return res
}

func (fbc *FileBlockChip) Decode(bb []byte) {
	var buf bytes.Buffer
	buf.Write(bb)
	dec := gob.NewDecoder(buf)
	dec.Decode(fbc)
}
