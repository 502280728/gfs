// 用户相关，用户是个接口
package gfs

import (
	"strconv"
)

type User interface {
	GetName() string        //获取用户名
	GetPass() string        //获取密码
	GetGroup() Group        //获取所在组
	GetUMask() string       //获取mask,三位，与Linux一致
	SetUMask(mask FileMask) //设置mask
}
type FileMask string

func (fm FileMask) GetAfterMasked() *FileAuth {
	fa := &FileAuth{}
	tmp := string(fm)
	fa.OwnerAuth = getMode(getUintMask(string(tmp[0])))
	fa.GroupAuth = getMode(getUintMask(string(tmp[1])))
	fa.OtherAuth = getMode(getUintMask(string(tmp[2])))
	return fa
}
func getUintMask(a string) uint32 {
	auth, _ := strconv.Atoi(a)
	return uint32(auth)
}
func getMode(a uint32) FileMode {
	fm := ALLMode
	m := FileMode(a)
	if m&ReadMode == ReadMode {
		fm ^= ReadMode
	}
	if m&WriteMode == WriteMode {
		fm ^= WriteMode
	}
	if m&ExecuteMode == ExecuteMode {
		fm ^= ExecuteMode
	}
	return fm
}

type Group string
