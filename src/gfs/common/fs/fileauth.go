package fs

import (
	"bytes"
)

//文件权限相关的东西，该权限系统与linux文件系统类似

type FileMode uint32

const (
	ReadMode FileMode = 1 << iota
	WriteMode
	ExecuteMode
	RWMode  FileMode = ReadMode | WriteMode
	REMode  FileMode = ReadMode | ExecuteMode
	WEMode  FileMode = WriteMode | ExecuteMode
	ALLMode FileMode = WriteMode | ExecuteMode | ReadMode
)

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

//新建文件夹时的默认权限
func NormalFileAuth() FileAuth {
	return FileAuth{OwnerAuth: ALLMode, GroupAuth: RWMode, OtherAuth: 0}
}
