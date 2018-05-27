package main

import (
	"fmt"
	"gfs/common"
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
	a := A{"Mike", 12}
	bb := common.EncodeToByteBuffer(a)

	tt := reflect.TypeOf(hello)
	rt := tt.In(0)

	aaa := reflect.New(rt)
	common.DecodeFromByteBuffer(aaa, bb)
	//	switch aaa.Interface().(type) {
	//	case *A:
	//		fmt.Println(aaa)
	//	}
	ccc, _ := aaa.Interface().(*A)
	fmt.Println(ccc)

	fu := reflect.ValueOf(hello)
	params := make([]reflect.Value, 1)
	params[0] = aaa.Elem()
	fu.Call(params)
}
