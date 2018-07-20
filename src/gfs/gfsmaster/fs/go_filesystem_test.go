// go_filesystem_test
package fs

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {

	myfs.MkDir("/a/b/c", true, rootUser)
	fmt.Println(myfs.List("/a", rootUser))
	fmt.Println("end")
}
