package fs

import (
	"gfs/gfsmaster/fs/user"
	"testing"
	"time"
)

func init() {
	Root = Node{Name: "", NodeFile: &File{Name: "", IsDir: true, CreateTime: time.Now()}, Nodes: []*Node{RootNode}}
}

var uuu = &user.User{Name: "Mike"}

func TestMakeDirAndList(t *testing.T) {
	fn := FileName("/ab/cd")

	fn.MakeDir(uuu)
	fn = FileName("/ab")
	f, _ := fn.List(uuu)
	if len(f) != 2 || f[1].Name != "/ab/cd" {
		t.FailNow()
	}
}

func TestFind(t *testing.T) {
	fn := FileName("/w/b.txt")
	fn.Touch(uuu)
	if f := fn.Find(uuu); f == nil || f.Name != "/w/b.txt" {
		t.FailNow()
	}
}

func TestDuplicateTouch(t *testing.T) {
	fn := FileName("/w/c.txt")
	_, _, e1 := fn.Touch(uuu)
	_, _, e2 := fn.Touch(uuu)
	if e1 != nil || e2 == nil || e2.Error() != "该文件已经存在" {
		t.FailNow()
	}
}
