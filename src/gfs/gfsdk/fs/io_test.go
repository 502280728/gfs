// io_test
package fs

import (
	"testing"
)

func TestWritor(t *testing.T) {
	writer := NewGFSWriter("/d.txt")
	bb := []byte("this is a test,and we will do nothing about it")
	writer.Write(bb)
	writer.Close()
}
