package fs

import (
	"fmt"
	"gfs/common"
	"gfs/common/gfs"
	"gfs/gfsmaster/wal"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GFileSystem struct {
	sync.Mutex
	TransactionId
	Root *GNode
}

type TransactionId struct {
	sync.Mutex
	Starts int64
}

func (ti *TransactionId) next() string {
	ti.Lock()
	defer ti.Unlock()
	ti.Starts++
	tmp := strconv.FormatInt(ti.Starts, 10)
	tmp = "00000000000000000000" + tmp
	return tmp[len(tmp)-20:]
}

func (ti *TransactionId) current() string {
	ti.Lock()
	defer ti.Unlock()
	tmp := strconv.FormatInt(ti.Starts, 10)
	tmp = "00000000000000000000" + tmp
	return tmp[len(tmp)-20:]
}

var myfs = getEmptyGFileSystem()

func getEmptyGFileSystem() *GFileSystem {
	root := &GNode{
		Name: "",
		NodeFile: &GFile{
			File: gfs.File{
				Name:       "",
				Owner:      rootUser,
				Auth:       gfs.DefaultFileAuth(),
				CreateTime: time.Now(),
				ModifyTime: time.Now(),
				Type:       gfs.TypeDirectory,
			},
			Visiable:  true,
			Locations: make([]*common.FileLocation, 0),
		},
		Nodes: make([]*GNode, 0, 3),
	}
	return &GFileSystem{
		Root: root,
	}
}

//为某一个文件分配一个block
func (fs *GFileSystem) Allocate(path string) {

}

func (fs *GFileSystem) replayWAL(w *wal.WAL, user gfs.User) {
	switch w.Op {
	case wal.OP_RM_FILE:
		fs.Remove(w.SourceFile, true, user)
	case wal.OP_CREATE_DIR:
		fs.MkDir(w.SourceFile, true, user)
	case wal.OP_CREATE_FILE:
		fs.Touch(w.SourceFile, user)
	}
}

func (fs *GFileSystem) MkDir(path string, recurisive bool, user gfs.User) (*gfs.File, error) {
	fs.Lock()
	defer fs.Unlock()
	paths := strings.Split(check(path), "/")
	index, gn := getLastGNodeAndIndex(paths)
	if index == len(path)-1 {
		return nil, fmt.Errorf(DIR_EXISTS_ERROR_FMT, path)
	}
	if recurisive || len(paths[(index+1):]) == 1 {
		node := createGNodes(gn, paths[(index+1):], recurisive, user, gfs.TypeDirectory)
		return &node.NodeFile.File, nil
	} else {
		return nil, fmt.Errorf(CREATE_DIR_ERROR_FMT, path)
	}
}
func (fs *GFileSystem) Touch(path string, user gfs.User) (*gfs.File, error) {
	fs.Lock()
	defer fs.Unlock()
	paths := strings.Split(check(path), "/")
	if gn := getGNode(paths[0:len(paths)-1], true); gn == nil {
		return nil, fmt.Errorf(CREATE_FILE_ERROR_FMT, path)
	} else {
		node := createGNodes(gn, paths[len(paths)-1:], true, user, gfs.TypeFile)
		return &node.NodeFile.File, nil
	}
}
func (fs *GFileSystem) Exists(path string, user gfs.User) (bool, error) {
	fs.Lock()
	defer fs.Unlock()
	paths := strings.Split(check(path), "/")
	if gn := getGNode(paths, false); gn == nil {
		return false, nil
	} else {
		return true, nil
	}
}
func (fs *GFileSystem) List(path string, user gfs.User) ([]*gfs.File, error) {
	fs.Lock()
	defer fs.Unlock()
	paths := strings.Split(check(path), "/")
	if gn := getGNode(paths, true); gn == nil {
		return nil, fmt.Errorf(DIR_NOT_EXISTS_ERROR_FMT, path)
	} else {
		res := make([]*gfs.File, 0, len(gn.Nodes))
		for _, n := range gn.Nodes {
			res = append(res, &n.NodeFile.File)
		}
		return res, nil
	}
}
func (fs *GFileSystem) GetFileInfo(path string, user gfs.User) (*gfs.File, error) {
	fs.Lock()
	defer fs.Unlock()
	paths := strings.Split(check(path), "/")
	if gn := getGNode(paths, false); gn != nil {
		return &gn.NodeFile.File, nil
	} else {
		return nil, fmt.Errorf(FILE_NOT_EXISTS_ERROR_FMT, path)
	}
}
func (fs *GFileSystem) Remove(path string, recurisive bool, user gfs.User) (*gfs.File, error) {
	fs.Lock()
	defer fs.Unlock()
	paths := strings.Split(check(path), "/")
	if gn := getGNode(paths, false); gn != nil {
		if gn.NodeFile.Type.IsFile() || gn.NodeFile.Type.IsLinkFile() || recurisive || (!recurisive && len(gn.Nodes) == 0) {
			par := getGNode(paths[:len(paths)-1], true)
			removeGNode(par, gn)
			return &gn.NodeFile.File, nil
		} else {
			return nil, fmt.Errorf(REMOVE_FILE_ERROR_FMT, path)
		}
	} else {
		return nil, fmt.Errorf(FILE_NOT_EXISTS_ERROR_FMT, path)
	}
}

func removeGNode(parent *GNode, target *GNode) {
	if len(parent.Nodes) == 1 {
		parent.Nodes = parent.Nodes[0:0]
		return
	}
	var tInd = -1
	for ind, node := range parent.Nodes {
		if node.NodeFile.Name == target.NodeFile.Name {
			tInd = ind
			break
		}
	}
	lastInd := len(parent.Nodes) - 1
	if tInd == lastInd {
		parent.Nodes = parent.Nodes[0:lastInd]
	} else {
		parent.Nodes[tInd] = parent.Nodes[lastInd]
		parent.Nodes = parent.Nodes[0:lastInd]
	}
}

//根据文件名在文件树中查找GNode。如果存在，则返回该GNode，如果不存在，返回nil
//isDir表示需要被查找的文件是文件夹还是文件。如果最后找到的GNode是文件，但是需要的是文件夹，仍然返回nil
func getGNode(paths []string, mustDir bool) *GNode {
	//如果paths长度为1，那么说明是根目录
	if len(paths) == 1 {
		return myfs.Root
	}
	res := myfs.Root
	for i := 1; i < len(paths); i++ {
		res = findIfNodeExists(paths[i], res.Nodes)
		if res == nil {
			break
		}
	}
	if mustDir && res != nil && res.NodeFile.Type != gfs.TypeDirectory && res.NodeFile.Type != gfs.TypeLinkDirectory {
		return nil
	} else {
		return res
	}
}

//找到paths中，对应的最后一个存在的GNode及该GNode的名字在paths中的索引
func getLastGNodeAndIndex(paths []string) (int, *GNode) {
	res := myfs.Root
	index := 0
	tmp := myfs.Root
	if len(paths) > 1 {
		for i := 1; i < len(paths); i++ {
			if res = findIfNodeExists(paths[i], tmp.Nodes); res == nil {
				res = tmp
				index = i
				break
			} else {
				tmp = res
			}
		}
		return index - 1, res
	} else {
		return 0, res
	}
}

//在某个GNode后面创建对应的文件夹,返回最后创建的文件夹对应的GNode
//paths为需要在p后面创建的文件夹
//recurisive表示是否需要递归创建
func createGNodes(p *GNode, paths []string, recurisive bool, user gfs.User, typ gfs.FileType) *GNode {
	if len(paths) > 1 && !recurisive {
		return nil
	}
	var tmp1, tmp2 *GNode
	tmp1 = p
	for _, path := range paths {
		tmp2 = &GNode{
			Name: path,
			NodeFile: &GFile{
				File: gfs.File{
					Name:       tmp1.NodeFile.Name + "/" + path,
					Owner:      user,
					Auth:       user.GetUMask().GetAfterMasked(),
					CreateTime: time.Now(),
					ModifyTime: time.Now(),
					Type:       typ,
				},
				Visiable:  true,
				Locations: make([]*common.FileLocation, 0),
			},
			Nodes: make([]*GNode, 0, 3),
		}
		tmp1.Nodes = append(tmp1.Nodes, tmp2)
		tmp1 = tmp2
	}
	return tmp2
}

func findIfNodeExists(name string, nodes []*GNode) *GNode {
	var res *GNode
	for _, node := range nodes {
		if name == node.Name {
			res = node
			break
		}
	}
	return res
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

type GNode struct {
	Name     string //节点名称，就是nodefile中文件的简称
	NodeFile *GFile
	Nodes    []*GNode
}

type GFile struct {
	gfs.File
	Visiable  bool //是否可见
	Locations []*common.FileLocation
}

type GUser struct {
	Name  string
	Pass  string
	Group gfs.Group
	Mask  gfs.FileMask
}

var rootUser = CreateUser("root", "root", gfs.Group("root"), gfs.DefaultFileMask)

func CreateUser(name, pass string, group gfs.Group, mask gfs.FileMask) gfs.User {
	return &GUser{name, pass, group, mask}
}

func (u *GUser) GetName() string {
	return u.Name
}
func (u *GUser) GetPass() string {
	return u.Pass
}
func (u *GUser) GetGroup() gfs.Group {
	return u.Group
}
func (u *GUser) GetUMask() gfs.FileMask {
	return u.Mask
}
