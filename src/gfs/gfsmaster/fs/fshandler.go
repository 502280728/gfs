package fs

import (
	"bytes"
	"gfs/common"
	"gfs/gfsmaster/fs/user"
)

func CreateHandler(string path, u *user.User, body string) []byte {
	fn := FileName(body)
	result := common.MessageInFS{}
	switch path {
	case "/fs/mkdir":
		succ, _, err := fn.MakeDir(u)
		result.Success = succ
		result.Msg = err.Error()
	case "/fs/touch":
		succ, _, err := fn.Touch(u)
		result.Success = succ
		result.Msg = err.Error()
	case "/fs/ls":
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

	}

}
