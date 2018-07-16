// go_filesystem_test
package fs

import (
	"fmt"
	"gfs/common/gfs"
	"testing"
)

func Test1(t *testing.T) {
	//	u := CreateUser("Mike", "123", "Mike", gfs.DefaultFileMask)
	//	fmt.Println(u.GetName())
	//	fmt.Println(u.GetUMask().GetAfterMasked())
	f := gfile{file: &gfs.File{Name: "mike"}, visiable: true}
	fmt.Println(f.file)
}
