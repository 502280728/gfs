// msg_test
package gmsg

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var a = MsgToMaster1{Block: 1, GFSMsg: GFSMsg{"2.0"}}
	fmt.Println(a)
	buf := a.Encode()
	var b MsgToMaster1
	b.Decode(buf.Bytes())
	fmt.Println(b.GFSMsg.Version)
	fmt.Println(b.Block)
}
