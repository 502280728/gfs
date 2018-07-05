// filesystem
package fs

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"
)

//表示一个文件或者文件夹
type File struct {
	Name       string    //文件的全名，末尾不包含"/"
	Owner      User      //所有者
	Mode       FileAuth  //权限
	CreateTime time.Time // 创建时间
	IsDir      bool      //是否是文件夹
	Size       int64     //文件大小，以字节计
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
	bb.WriteString(strconv.FormatInt(file.Size, 10))
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

//表示一个文件系统
type FileSystem interface {
	Create(path string) (io.WriteCloser, error)         //创建一个文件
	Open(path string) (io.ReadCloser, error)            //打开一个已经存在的文件
	MkDir(path string) (*File, error)                   //新建文件夹
	Touch(path string) (*File, error)                   //新建一个文件
	Exists(paht string) (bool, error)                   //判断一个文件或者文件夹是否存在
	List(path string) ([]*File, error)                  //获取文件夹的子文件夹/文件夹
	GetFileInfo(path string) (*File, error)             //获取文件/文件夹详情
	Remove(path string, recurisive bool) (*File, error) //若删除成功，error为nil，否则不为nil；FileInfo为删除文件的信息
}
