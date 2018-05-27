package common

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func Test(t *testing.T) {
	p := map[string]string{"a": "b"}
	bb := EncodeToByteBuffer(p)
	fmt.Println("decode")
	pppp := map[string]string{}
	DecodeFromByteBuffer(&pppp, bb)
	fmt.Println(pppp)
}
