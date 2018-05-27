package fs

import (
	"bytes"
	"gfs/common"
	"gfs/gfsmaster/fs/user"
)

func Handle(path string, u *user.User, body string) []byte {
	fn := FileName(body)
	result := common.MessageInFS{}
	switch path {

	case "/fs/mkdir":
		succ, _, err := fn.MakeDir(u)
		result.Success = succ
		if !succ {
			result.Msg = err.Error()
		}
	case "/fs/touch":
		succ, _, err := fn.Touch(u)
		result.Success = succ
		if !succ {
			result.Msg = err.Error()
		}
	case "/fs/ls", "/fs/ll":
		if files, err := fn.List(u); err == nil {
			result.Success = true
			var bb bytes.Buffer
			for _, file := range files {
				bb.WriteString(file.String() + "\n")
			}
			result.Data = string(bb.Bytes())
		} else {
			result.Success = false
			result.Msg = err.Error()
		}
	case "/fs/rm":
		if _, err := fn.Remove(u); err == nil {
			result.Success = true
			result.Data = ""
		} else {
			result.Success = false
			result.Msg = err.Error()
		}
	case "/fs/chmod":
	case "/fs/chown":
	case "/fs/adduser":
	}
	return common.EncodeToBytes(result)
}
