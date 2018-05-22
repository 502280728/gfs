package main

import (
	"fmt"
	"gfs/common"
)

func main() {
	a := &common.GFSReader{TargetFile: "/wzf/a.txt", TargetNode: []*common.FileLocation{&common.FileLocation{Main: "localhost:8087"}}}
	pp := make([]byte, 1024, 1024)
	a.Read(pp)
	fmt.Println(string(pp))
	//	var fbc = &common.FileBlockChip{FileName: "/wzf/a.txt", Limit: 10, Offset: 0, Block: 0}
	//	buf := fbc.Encode()
	//	bb, _ := common.GetDataFromSpecialURL(buf, "http://localhost:8087")
	//	fmt.Println(string(bb))
}
