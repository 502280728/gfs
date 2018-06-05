package main

import (
	"fmt"
	"gfs/common/http"
	"reflect"
)

type A struct {
	Name string
	Age  int
}

func (a A) String() string {
	return fmt.Sprintf("%s is %d years old。", a.Name, a.Age)
}

func hello(a A) {
	fmt.Println(a.Name)
}

func main() {

}
