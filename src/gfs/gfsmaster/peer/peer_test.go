// peer_test
package peer

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	var p Peer
	p.LoadConf("d:/temp/conf/gfs.properties")
	fmt.Println(p.Conf.Address.FriendIP, "123")
	var a string
	fmt.Println(p.Conf.Address.FriendIP == "")
	fmt.Println(a == "")
}
