// peer_test
package peer

import (
	"testing"
)

func Test1(t *testing.T) {
	InitPeer("d:/temp/conf/gfs.properties")
	GetPeer()
}
