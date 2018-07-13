// 由sdk往外发出的信息
package gmsg

import (
	"bytes"
	"gfs/common/gutils"
)

//sdk发往master获取每一个block存储位置的消息
type MsgToMaster1 struct {
	Block    int
	FileName string
}

func (msg *MsgToMaster1) Encode() *bytes.Buffer {
	return gutils.EncodeToByteBuffer(msg)
}

func (msg *MsgToMaster1) Decode(bb []byte) {
	gutils.DecodeFromBytes(msg, bb)
}
