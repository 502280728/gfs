package main

import (
	"fmt"
	"os"

	"reflect"

	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gfs/test")

type A struct {
	Name string
	Age  int
}

func (a A) String() string {
	return fmt.Sprintf("%s is %d years old。", a.Name, a.Age)
}

func hello(a A) {

}

func main() {
	//	a := map[string]string{"name": "Mike"}
	//	bb := common.EncodeToByteBuffer(a)
	//	res := map[string]string{}
	//	common.DecodeFromByteBuffer(&res, bb)
	//	fmt.Println(res)
	//	var a fs.FileName = "/a/b"
	//	b := &a
	//	logger.Infof("%s", string(*b))
	//	s, _ := os.Stat("d:/temp/a.txt")
	//	fmt.Print(s.Size())
	//	file, _ := os.OpenFile("d:/temp/a.txt", os.O_APPEND, os.ModeAppend)
	//	defer file.Close()
	//	for i := 0; i < 180; i++ {
	//		file.WriteString("abcdefgh" + strconv.Itoa(i) + "\r\n")
	//	}
	//	l, _ = file.Read(bb)
	//	fmt.Println(l, string(bb))
	//ioutil.WriteFile("d:\temp\ab.txt", []byte("abc"), os.ModeAppend)
	os.Create("d:/temp/data/abc.txt")
	ff, _ := os.OpenFile("d:/temp/data/abc.txt", os.O_APPEND, os.ModeAppend)
	defer ff.Close()
	ff.Write([]byte("abcc"))
	fmt.Println(reflect.Value{} != reflect.Value{})
}
