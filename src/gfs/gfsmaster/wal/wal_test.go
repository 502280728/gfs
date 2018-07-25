// wal_test
package wal

import (
	"fmt"
	"gfs/gfsmaster/common"
	"testing"
)

func Test1(t *testing.T) {
	mcommon.LoadConf("d:/temp/conf/gfs.properties")
	walConf = mcommon.GetPeerConf().WAL
	var abc WALFileReader
	buf := make([]*WAL, 2, 2)
	ind, _ := abc.Read("00000000000000000001", buf)
	fmt.Println(ind)
	fmt.Println(string(buf[0].ToBytes()))
	fmt.Println(string(buf[1].ToBytes()))
	ind, _ = abc.Read("00000000000000000003", buf)
	fmt.Println(ind)
	fmt.Println(string(buf[0].ToBytes()))
	fmt.Println(string(buf[1].ToBytes()))
}
