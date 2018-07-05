// go_filesystem
package fs

import (
	"fmt"
	"gfs/common/fs"
	"io"
)

type GoFileSystem struct {
}

func (gfs *GoFileSystem) Open(path string) (io.ReadCloser, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Create(path string) (io.WriteCloser, error) {
	if found, err := gfs.Exists(path); found || err != nil {
		if found {
			return nil, fmt.Errorf("文件%s已经存在", path)
		} else {
			return nil, err
		}
	}

	return nil, nil
}

func (gfs *GoFileSystem) MkDir(path string) (*fs.File, error) {
	return nil, nil
}

func (gfs *GoFileSystem) Exists(path string) (bool, error) {
	return true, nil
}
