package main

import (
	"flag"
	"fmt"
	"strings"
)

//文件系统的操作命令,目前支持 rm,copy,mkdir
// cp -r source target
// cp source target
// rm -r target
// mkdir -p target
// mkdir target
func main() {
	k := strings.Split("   -r=l", " ")
	w := flag.String("r", "a", "a string")
	flag.CommandLine.Parse(k)
	fmt.Println(*w)
}
