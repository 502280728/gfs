package fs

import (
	"bytes"
	"gfs/gfsmaster/fs/user"
	"log"
	"strings"
	"time"
)

type FileSystem interface {
	Mkdir(name string)
	Remove(name string, recursive bool)
	Touch(name string)
}

//代表一个文件或者文件夹
//有如下的规定：根目录"/"的Name和ParentDir都为""
type File struct {
	Name       string    //文件名
	ParentDir  string    //父文件夹名称,在实际存储中，该值并不存在，为空仅当新建时才行;该dir必须以“/”开头,结尾不包含"/"
	Owner      user.User //所有者
	Mode       FileAuth  //权限
	CreateTime time.Time // 创建时间
	IsDir      bool
}

//文件系统在内存中以树的形式存在。
//这是树的节点.
type Node struct {
	NodeFile *File   //当前节点的文件
	Nodes    []*Node //子节点
}

var Root = Node{NodeFile: &File{Name: "", ParentDir: ""}, Nodes: []*Node{}}

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
	bb.WriteString(file.Name)
	return string(bb.Bytes())
}

func (file *File) CanRead(user *user.User) bool {
	return can(file, user, ReadMode)
}
func (file *File) CanWrite(user *user.User) bool {
	return can(file, user, ReadMode)
}
func (file *File) CanExecute(user *user.User) bool {
	return can(file, user, ExecuteMode)
}

func can(file *File, user *user.User, fm FileMode) bool {
	if user.Name == file.Owner.Name {
		return checkMode(file.Mode.OwnerAuth, fm)
	} else if user.Group == file.Owner.Group {
		return checkMode(file.Mode.GroupAuth, fm)
	} else {
		return checkMode(file.Mode.OtherAuth, fm)
	}
}

func checkMode(fm FileMode, target FileMode) bool {
	return fm|target == target
}

func List(file *File) []*File {
	fs := []*File{}
	for _, node := range searchFromRoot(file.ParentDir + "/" + file.Name).Nodes {
		fs = append(fs, node.NodeFile)
	}
	return fs
}

//逻辑是先把file的ParentDir建立或者找到
func Create(file *File) (bool, error) {
	log.Printf("creating file %s/%s", file.ParentDir, file.Name)
	names := splitPath(file.ParentDir)
	var node *Node = &Root
	var ind int
	for index, name := range names {
		var find = false
		for _, tmp := range node.Nodes {
			if name == tmp.NodeFile.Name {
				node = tmp
				find = true
				break
			}
		}
		ind = index
		if !find {
			break
		}
	}

	var begin, end *Node
	if ind < len(names) {
		begin, end = createNodes(names[ind:], file)
		node.Nodes = append(node.Nodes, begin)
	} else {
		end = node
	}
	end.Nodes = append(end.Nodes, &Node{NodeFile: &File{Name: file.Name, Owner: file.Owner, Mode: file.Mode, CreateTime: time.Now(), IsDir: false}})
	return true, nil
}

//建立parentdir所有的node，返回新建的第一个和最后一个node
func createNodes(names []string, file *File) (*Node, *Node) {
	node := &Node{}
	if len(names) == 1 {
		node.NodeFile = &File{Name: file.Name, Owner: file.Owner, Mode: file.Mode, CreateTime: time.Now(), IsDir: true}
		return node, node
	} else {
		node = &Node{NodeFile: &File{Name: names[0], Owner: file.Owner, Mode: file.Mode, CreateTime: time.Now(), IsDir: true}, Nodes: []*Node{}}
		tmp := node
		for _, name := range names[1:] {
			tmpNode := &Node{NodeFile: &File{Name: name, Owner: file.Owner, Mode: file.Mode, CreateTime: time.Now(), IsDir: true}, Nodes: []*Node{}}
			tmp.Nodes = append(tmp.Nodes, tmpNode)
			tmp = tmpNode
		}
		return node, tmp
	}
}

//找到文件或者文件夹对应的Node节点
func searchFromRoot(path string) *Node {
	if path == "/" {
		return &Root
	}
	names := splitPath(path)
	var node *Node = &Root
	for _, name := range names {
		var find = false
		for _, tmp := range node.Nodes {
			if name == tmp.NodeFile.Name {
				node = tmp
				find = true
				break
			}
		}
		if !find {
			return nil
		}
	}
	return node
}

func splitPath(path string) []string {
	return strings.Split(path, "/")[1:]
}
