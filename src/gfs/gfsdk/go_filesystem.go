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

func (gfs *GoFileSystem) Open(path string, user User) (gio.ReadCloser, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Create(path string, user User) (gio.WriteCloser, error) {
	if found, err := gfs.Exists(path); found || err != nil {
		if found {
			return nil, fmt.Errorf("文件%s已经存在", path)
		} else {
			return nil, err
		}
	}

	return NewGFSWriter(path), nil
}

func (gfs *GoFileSystem) MkDir(path string, user User) (*gfs.File, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Touch(path string, user User) (*gfs.File, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Exists(path string, user User) (bool, error) {
	return false, nil
}
func (gfs *GoFileSystem) List(path string, user User) ([]*gfs.File, error) {
	return nil, nil
}
func (gfs *GoFileSystem) GetFileInfo(path string, user User) (*gfs.File, error) {
	return nil, nil
}
func (gfs *GoFileSystem) Remove(path string, recurisive bool, user User) (*gfs.File, error) {
	return nil, nil
}
