package fs

import (
	"gfs/common"
	"gfs/common/gfs"
	"gfs/common/gio"
	"strings"
)

type gfilesystem struct {
	root *gnode
}

func (fs *gfilesystem) Create(path string, user gfs.User) (gio.WriteCloser, error) {
	return nil, nil
}
func (fs *gfilesystem) Open(path string, user gfs.User) (gio.ReadCloser, error) {
	return nil, nil
}
func (fs *gfilesystem) MkDir(path string, user gfs.User) (*gfs.File, error) {

	return nil, nil
}
func (fs *gfilesystem) Touch(path string, user gfs.User) (*gfs.File, error) {
	return nil, nil
}
func (fs *gfilesystem) Exists(paht string, user gfs.User) (bool, error) {
	return false, nil
}
func (fs *gfilesystem) List(path string, user gfs.User) ([]*gfs.File, error) {
	return nil, nil
}
func (fs *gfilesystem) GetFileInfo(path string, user gfs.User) (*gfs.File, error) {
	return nil, nil
}
func (fs *gfilesystem) Remove(path string, recurisive bool, user gfs.User) (*gfs.File, error) {
	return nil, nil
}

//验证文件名是否合法,同时将一些偏门写法转换为正确的写法：
// /a/b/c/   -> /a/b/c
// /a/b//c/  -> /a/b/c
// \a\b\c\   -> /a/b/c
// \a\b\c    -> /a/b/c
// /a\b/c    -> /a/b/c
// TODO
func check(path string) string {
	return strings.TrimSuffix(path, "/")
}

type gnode struct {
	name     string //节点名称，就是nodefile中文件的简称
	nodeFile gfile
	nodes    []*gnode
}

type gfile struct {
	file      *gfs.File
	visiable  bool //是否可见
	locations []*common.FileLocation
}

type guser struct {
	name  string
	pass  string
	group gfs.Group
	mask  gfs.FileMask
}

func CreateUser(name, pass string, group gfs.Group, mask gfs.FileMask) gfs.User {
	return &guser{name, pass, group, mask}
}

func (u *guser) GetName() string {
	return u.name
}
func (u *guser) GetPass() string {
	return u.pass
}
func (u *guser) GetGroup() gfs.Group {
	return u.group
}
func (u *guser) GetUMask() gfs.FileMask {
	return u.mask
}
