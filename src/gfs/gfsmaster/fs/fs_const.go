// fs_const
package fs

import (
	"errors"
)

var (
	DIR_EXISTS_ERROR      = errors.New("该文件夹已经存在")
	FILE_EXISTS_ERROR     = errors.New("该文件已经存在")
	DIR_NOT_EXISTS_ERROR  = errors.New("该文件夹不存在")
	FILE_NOT_EXISTS_ERROR = errors.New("该文件不存在")
	CREATE_DIR_ERROR      = errors.New("创建文件夹失败")
)
