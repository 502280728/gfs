package main

import (
	"gfs/gfsmaster/fs"
	"gfs/gfsmaster/fs/user"
)

func main() {

	file := &fs.File{Name: "test1.txt", ParentDir: "/a/b/c", Owner: user.User{Name: "root", Group: "root"}, Mode: fs.NormalFileAuth(), IsDir: false}

	fs.Create(file)
	file = &fs.File{Name: "test2.txt", ParentDir: "/a/b", Owner: user.User{Name: "root", Group: "root"}, Mode: fs.NormalFileAuth(), IsDir: false}
	fs.Create(file)
	file = &fs.File{Name: "test2.txt", ParentDir: "", Owner: user.User{Name: "root", Group: "root"}, Mode: fs.NormalFileAuth(), IsDir: false}
	fs.Create(file)
	println("begin")

	for _, file := range fs.List(&fs.File{Name: "", ParentDir: ""}) {
		println(file.String())
	}

}
