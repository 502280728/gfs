// context_test
package gfsdk

import (
	"testing"
)

func TestFileSystem(t *testing.T) {
	context := GetContext()
	context.LoadConf("D:/temp/conf/gfs.properties")
	fs := context.GetDefaultFileSystem()
	bb := []byte("this is a test,and we will do nothing about itthis is a test,and we will do nothing about itthis is a test,and we will do nothing about it")
	writer, _ := fs.Create("d:/a.txt")
	writer.Write(bb)
	writer.Close()
}
