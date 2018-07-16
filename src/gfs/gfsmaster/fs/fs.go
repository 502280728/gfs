//可以认为是common.gfs的实现
//包含了文件系统和用户系统两个子系统
//其中文件系统的实现方式是文件树
package fs

import (
	"gfs/common/gfs"
)

type GFS struct {
	fs    gfs.FileSystem
	users []gfs.User
}
