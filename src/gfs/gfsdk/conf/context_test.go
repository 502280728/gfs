// context_test
package conf

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	GFSContext.LoadConf("D:/temp/conf/gfs.properties")
	GFSContext.GetFileSystem()

	conf := GFSContext.GetConf()
	fmt.Println(conf.Master)
	fmt.Println(GFSContext.GetConf().BlockSize)
}
