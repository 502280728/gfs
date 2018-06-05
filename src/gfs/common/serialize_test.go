package common

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func Test(t *testing.T) {
	p := &Person{}
	pv := reflect.ValueOf(p).Elem()
	pv.FieldByName("Name").Set(reflect.ValueOf("mike"))
	fmt.Println(p)
}
