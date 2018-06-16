package main

import (
	"fmt"
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

}

func main() {
	fmt.Println(reflect.Value{} != reflect.Value{})
}
