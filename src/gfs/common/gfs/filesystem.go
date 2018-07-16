// 一个文件系统包括了两个子系统：
//文件和用户
package gfs

import (
	"bytes"
	"gfs/common/gio"
	"strconv"
	"strings"
	"time"
)

//表示一个文件或者文件夹
type File struct {
	Name          string    //文件的全名，末尾不包含"/"
	Owner         User      //所有者
	Mode          FileAuth  //权限
	CreateTime    time.Time // 创建时间
	ModifyTime    time.Time //修改时间
	LastVisitTime time.Time //最后一个访问时间
	Type          FileType  // 文件类型
	Size          int64     //文件大小，以字节计
	LinkedFile    *File     //如果文件类型是软连接，那么这里存放目标文件的指针
}

//表示一个文件系统
type FileSystem interface {
	Create(path string) (gio.WriteCloser, error)        //创建一个文件
	Open(path string) (gio.ReadCloser, error)           //打开一个已经存在的文件
	MkDir(path string) (*File, error)                   //新建文件夹
	Touch(path string) (*File, error)                   //新建一个文件
	Exists(paht string) (bool, error)                   //判断一个文件或者文件夹是否存在
	List(path string) ([]*File, error)                  //获取文件夹的子文件夹/文件夹
	GetFileInfo(path string) (*File, error)             //获取文件/文件夹详情
	Remove(path string, recurisive bool) (*File, error) //若删除成功，error为nil，否则不为nil；FileInfo为删除文件的信息
}

func (file *File) String() string {
	bb := bytes.Buffer{}
	if file.Type.IsDir() {
		bb.WriteString("d")
	} else if file.Type.IsFile() {
		bb.WriteString("-")
	} else if file.Type.IsLink() {
		bb.WriteString("l")
	} else {
		bb.WriteString(" ")
	}
	bb.WriteString(file.Mode.String())
	bb.WriteByte('\t')
	bb.WriteString(file.Owner.GetName())
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
