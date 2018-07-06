// go_filesystem
package gfsdk

import (
	"fmt"
	"gfs/common/gfs"
	"gfs/common/gio"
)

type GoFileSystem struct {
	Name string
}

func (gfs *GoFileSystem) Open(path string) (gio.ReadCloser, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Create(path string) (gio.WriteCloser, error) {
	if found, err := gfs.Exists(path); found || err != nil {
		if found {
			return nil, fmt.Errorf("文件%s已经存在", path)
		} else {
			return nil, err
		}
	}

	return NewGFSWriter(path), nil
}

func (gfs *GoFileSystem) MkDir(path string) (*gfs.File, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Touch(path string) (*gfs.File, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Exists(path string) (bool, error) {
	return false, nil
}
func (gfs *GoFileSystem) List(path string) ([]*gfs.File, error) {
	return nil, nil
}
func (gfs *GoFileSystem) GetFileInfo(path string) (*gfs.File, error) {
	return nil, nil
}
func (gfs *GoFileSystem) Remove(path string, recurisive bool) (*gfs.File, error) {
	return nil, nil
}
