// GFS可能需要特制的io，所以这里先仿照标准库建立一个io，
package gio

import (
	"errors"
)

var EOF = errors.New("End OF File")

type ReadCloser interface {
	Read(p []byte) (int, error)
	Close() error
}

type WriteCloser interface {
	Write(p []byte) (int, error)
	Close() error
	Flush() error
}
