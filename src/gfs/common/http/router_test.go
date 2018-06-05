// router_test.go
package http

import (
	"fmt"
	"gfs/common"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func HH(p Person) Person {
	return p
}

func Test1(t *testing.T) {
	person := Person{"Mike", 12}
	bb := common.EncodeToByteBuffer(person)
	tt := handle(HH, bb, GO)

	cc := common.EncodeToByteBuffer(tt)
	pp := &Person{}
	common.DecodeFromByteBuffer(pp, cc)
	fmt.Println(pp)
}
