package fs

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"gfs/gfsmaster/fs/user"
	logging "github.com/op/go-logging"
	"os"
	"strings"
	"time"
)

var logger = logging.MustGetLogger("gfs/gfsmaster/fs")

//代表一个文件或者文件夹
//有如下的规定：根目录"/"的Name和ParentDir都为""
type File struct {
	Name       string    //文件的全名，末尾不包含"/"
	Owner      user.User //所有者
	Mode       FileAuth  //权限
	CreateTime time.Time // 创建时间
	IsDir      bool      //是否是
	Unvisiable bool      //是否可见
	Paths      []string  //在物理节点上的存储位置，格式是 [ip:port;id,........]
}

//有些场景仅仅需要传递文件名称，该type包装这个名称
type FileName string //支持 /a ,/a/,/a/b.txt三种形式

//                            "/"
//   "a" {"a",File{"/a",....}}    "b"
//

//文件系统在内存中以树的形式存在。
//这是树的节点.
type Node struct {
	Name     string  //节点名称，就是NodeFile中的文件的简称
	NodeFile *File   //当前节点的文件
	Nodes    []*Node //子节点
}

//目录的根节点
var RootNode = &Node{NodeFile: &File{Name: ".", IsDir: true, CreateTime: time.Now(), Mode: NormalFileAuth()}}
var Root = Node{}

var RootStorePath string

func init() {
	logger.Info("starting filesystem")
}

func RecoverFromStore() error {
	if _, err := os.Stat(RootStorePath); err == nil {
		logger.Infof("recover file system from file: %s \r\n", RootStorePath)
		if file, err := os.Open(RootStorePath); err == nil {
			dec := gob.NewDecoder(file)
			err2 := dec.Decode(&Root)
			return err2
		} else {
			return err
		}
	} else {
		logger.Warningf("the file %s does not exist, start a empty filesystem", RootStorePath)
		Root = Node{Name: "", NodeFile: &File{Name: "", IsDir: true, CreateTime: time.Now()}, Nodes: []*Node{RootNode}}
		return nil
	}
}

func StoreFileSystem() error {
	logger.Infof("store file system in file: %s \r\n", RootStorePath)
	if file, err := os.Create(RootStorePath); err == nil {
		enc := gob.NewEncoder(file)
		return enc.Encode(Root)
	} else {
		return err
	}

}

func (file *File) String() string {
	bb := bytes.Buffer{}
	if file.IsDir {
		bb.WriteString("d")
	} else {
		bb.WriteString("-")
	}
	bb.WriteString(file.Mode.String())
	bb.WriteByte('\t')
	bb.WriteString(file.Owner.Name)
	bb.WriteByte('\t')
	bb.WriteString(file.CreateTime.Format("2006-01-02 15:04:05"))
	bb.WriteByte('\t')
	bb.WriteString(simplifyName(file.Name))
	return string(bb.Bytes())
}

func simplifyName(name string) string {
	if name == "." || name == ".." {
		return name
	} else {

		return string([]byte(name)[strings.LastIndex(name, "/")+1:])
	}
}

func (node *Node) String() string {
	return node.NodeFile.Name
}

//验证文件名是否合法,同时将一些偏门写法转换为正确的写法：
// /a/b/c/   -> /a/b/c
// /a/b//c/  -> /a/b/c
// \a\b\c\   -> /a/b/c
// \a\b\c    -> /a/b/c
// /a\b/c    -> /a/b/c
// TODO
func (fn *FileName) check() (string, error) {
	return strings.TrimSuffix(string(*fn), "/"), nil
}

//新建文件夹时验证权限
func check(node *Node, user *user.User, mode FileMode) bool {
	return true
}

func (fn *FileName) List(user *user.User) ([]*File, error) {
	if name, err := fn.check(); err == nil {
		names := strings.Split(name, "/")
		if index, node := findNotExists(names); index == -1 {
			if check(node, user, ReadMode) {
				var result []*File
				for _, nod := range node.Nodes {
					if !nod.NodeFile.Unvisiable {
						result = append(result, nod.NodeFile)
					}
				}
				return result, nil
			} else {
				return nil, errors.New("权限不足")
			}
		} else {
			return nil, errors.New("文件或文件夹不存在")
		}

	} else {
		return nil, err
	}
}

//新建文件夹
//如果出现权限错误，或者文件名错误，返回fasle,nil,error
//如果文件夹已经存在，返回false，和对应的Node节点,
//如果文件夹不存在，返回true，和新建的最后的Node节点
func (fn *FileName) MakeDir(user *user.User) (bool, *Node, error) {
	if name, err := fn.check(); err == nil {
		names := strings.Split(name, "/")
		if index, node := findNotExists(names); index == -1 {
			return false, node, errors.New("文件夹已经存在") //已经存在的不检查权限
		} else {
			if check(node, user, WEMode) {
				begin, end := createNode(node.NodeFile.Name, names[index:], user)
				node.Nodes = append(node.Nodes, begin)
				return true, end, nil
			} else {
				return false, nil, errors.New("权限不足")
			}
		}
	} else {
		return false, nil, err
	}
}

func (fn *FileName) Touch(user *user.User) (bool, *Node, error) {
	nn := string(*fn)
	parentDir := FileName(string([]byte(nn)[0:strings.LastIndex(nn, "/")]))
	simpleName := string([]byte(nn)[strings.LastIndex(nn, "/")+1:])
	if _, node, err := parentDir.MakeDir(user); node != nil {
		if node.NodeFile.IsDir {
			if check(node, user, WriteMode) {
				newFile := &File{Name: string(*fn), IsDir: false, Mode: NormalFileAuth(), Owner: *user, CreateTime: time.Now()}
				newNode := &Node{Name: simpleName, NodeFile: newFile}
				node.Nodes = append(node.Nodes, newNode)
				return true, newNode, nil
			} else {
				return false, nil, errors.New("权限不足")
			}
		} else {
			return false, nil, fmt.Errorf("文件%s创建失败", *fn)
		}
	} else {
		return false, nil, err
	}

}

func (fn *FileName) Remove(user *user.User) (bool, error) {
	return false, nil
}

func createNode(parentDir string, names []string, user *user.User) (*Node, *Node) {
	var nodes []*Node
	for index, name := range names {
		file := &File{Name: parentDir + "/" + strings.Join(names[0:index+1], "/"), IsDir: true, Mode: NormalFileAuth(), Owner: *user, CreateTime: time.Now()}
		selfNode := &Node{NodeFile: &File{Name: ".", IsDir: true, Mode: NormalFileAuth(), Owner: *user, CreateTime: time.Now()}}
		nodes = append(nodes, &Node{Name: name, NodeFile: file, Nodes: []*Node{selfNode}})
	}
	if len(nodes) > 1 {
		for i := 0; i < len(nodes)-1; i++ {
			nodes[i].Nodes = append(nodes[i].Nodes, nodes[i+1])
		}
	}
	return nodes[0], nodes[len(nodes)-1]
}

//根据文件名在文件树中寻找最后一个不存在的文件夹，如果都存在，则返回最后一个文件夹所在的节点
func findNotExists(names []string) (int, *Node) {
	var resNode *Node = &Root
	find := false
	var tmp *Node
	for index, _ := range names {
		if index < len(names)-1 {
			if find, tmp = findExistsNode(names[index+1], resNode.Nodes); !find {
				return index + 1, resNode
			} else {
				resNode = tmp
			}
		}
	}
	return -1, resNode
}

func findExistsNode(name string, nodes []*Node) (bool, *Node) {
	var find = false
	var res *Node = nil
	for _, node := range nodes {
		if name == node.Name {
			find = true
			res = node
			break
		}
	}
	return find, res
}

func (file *File) CanBeReadBy(user *user.User) bool {
	return can(file, user, ReadMode)
}
func (file *File) CanBeWroteBy(user *user.User) bool {
	return can(file, user, ReadMode)
}
func (file *File) CanByExecutedBy(user *user.User) bool {
	return can(file, user, ExecuteMode)
}

func can(file *File, user *user.User, fm FileMode) bool {
	if user.Name == file.Owner.Name {
		return checkMode(file.Mode.OwnerAuth, fm)
	} else {
		return checkMode(file.Mode.OtherAuth, fm)
	}
}
func checkMode(fm FileMode, target FileMode) bool {
	return fm|target == target
}
