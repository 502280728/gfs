package main

import (
	"gfs/gfsmaster/server"
	"github.com/spf13/cobra"
)

var gfsmasterCmd = &cobra.Command{
	Use: "gfsmaster",
}

func main() {
	gfsmasterCmd.AddCommand(server.Cmd())
	gfsmasterCmd.Execute()

	//	fmt.Println(strings.Join([]string{"er"}[0:0], "/"))
	//
	//	user1 := &user.User{Name: "root", Group: "root"}
	//	//	dirs := []string{"/d/e", "/d/er"}
	//	//	for _, dir := range dirs {
	//	//		tmp := fs.FileName(dir)
	//	//		tmp.MakeDir(user1)
	//	//	}
	//	//	user1 = &user.User{Name: "root1", Group: "root1"}
	//	//	dirs = []string{"/d/e/a", "/d/b"}
	//	//	for _, dir := range dirs {
	//	//		tmp := fs.FileName(dir)
	//	//		tmp.MakeDir(user1)
	//	//	}
	//	//	fs.StoreFileSystem()
	//	fs.RootStorePath = "D:/temp/fs1.binary"
	//	fs.RecoverFromStore()
	//	fn1 := fs.FileName("/a.txt")
	//	fn1.Touch(user1)
	//	fn := fs.FileName("/")
	//
	//	n, _ := fn.List(user1)
	//
	//	for _, tmp := range n {
	//		fmt.Println(tmp)
	//	}
}
