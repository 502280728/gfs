//该文件包括文件系统的其他的常量或者数据结构
//包含了
//	文件类型：文件、文件夹、链接
//	文件权限
//	权限mask
package gfs

import (
	"bytes"
)

//文件类型，支持文件、文件夹、文件软连接、文件夹软连接
type FileType uint32

const (
	TypeDirectory     FileType = 1 << iota                //文件夹
	TypeFile                                              //文件
	TypeLink                                              //软连接
	TypeLinkFile      FileType = TypeLink | TypeFile      //文件软连接
	TypeLinkDirectory FileType = TypeLink | TypeDirectory //文件夹软连接
)

//文件权限相关的东西，该权限系统与linux文件系统类似
type FileMode uint32

const (
	ReadMode FileMode = 1 << iota
	WriteMode
	ExecuteMode
	RWMode   FileMode = ReadMode | WriteMode
	REMode   FileMode = ReadMode | ExecuteMode
	WEMode   FileMode = WriteMode | ExecuteMode
	ALLMode  FileMode = WriteMode | ExecuteMode | ReadMode
	NoneMode FileMode = 0
)

//是否是文件夹
func (ft FileType) IsDir() bool {
	return ft == TypeDirectory
}

//是否是文件
func (ft FileType) IsFile() bool {
	return ft == TypeFile
}

//是否是软连接
func (ft FileType) IsLink() bool {
	return ft&TypeLink == TypeLink
}

//是否是文件夹软连接
func (ft FileType) IsLinkDir() bool {
	return ft == TypeLinkDirectory
}

//是否会死文件软连接
func (ft FileType) IsLinkFile() bool {
	return ft == TypeLinkFile
}

func (fm FileMode) String() string {
	buf := []byte("---")
	if fm&ReadMode == ReadMode {
		buf[0] = 'r'
	}
	if fm&WriteMode == WriteMode {
		buf[1] = 'w'
	}
	if fm&ExecuteMode == ExecuteMode {
		buf[2] = 'x'
	}
	return string(buf)
}

//文件全部的权限，包括所有者权限、组权限、其他人员权限
type FileAuth struct {
	OwnerAuth FileMode
	GroupAuth FileMode
	OtherAuth FileMode
}

func (fa *FileAuth) String() string {
	bb := bytes.Buffer{}
	bb.WriteString(fa.OwnerAuth.String())
	bb.WriteString(fa.GroupAuth.String())
	bb.WriteString(fa.OtherAuth.String())
	return string(bb.Bytes())
}
