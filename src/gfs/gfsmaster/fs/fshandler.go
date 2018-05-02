package fs

import (
	"bytes"
	"encoding/gob"
	"gfs/common"
	"gfs/gfsmaster/fs/user"
)

func CreateHandler(path string, u *user.User, body string) []byte {
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
	}
	var bb1 bytes.Buffer
	enc := gob.NewEncoder(&bb1)
	enc.Encode(result)
	return bb1.Bytes()
}
