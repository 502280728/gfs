// go_filesystem_test
package fs

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	myfs.MkDir("/a/b/c", true, rootUser)
	myfs.MkDir("/a/c", true, rootUser)
	if a, err1 := myfs.List("/", rootUser); err1 != nil || len(a) != 1 || a[0].Name != "/a" {
		t.FailNow()
	}
	if b, err2 := myfs.List("/a", rootUser); err2 != nil || len(b) != 2 || b[0].Name != "/a/b" {
		t.FailNow()
	}

}

func Test2(t *testing.T) {
	myfs.MkDir("/rm", true, rootUser)
	myfs.MkDir("/rm/ab", true, rootUser)
	fmt.Println(myfs.List("/rm", rootUser))
	myfs.Remove("/rm/ab", true, rootUser)
	fmt.Println(myfs.List("/rm", rootUser))
	myfs.Touch("/rm/a.txt", rootUser)
	fmt.Println(myfs.List("/rm", rootUser))
}
