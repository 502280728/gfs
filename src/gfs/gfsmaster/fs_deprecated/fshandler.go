package fs_d

import (
	"bytes"
	"errors"
	"gfs/common"
	"gfs/common/utils"
	"gfs/gfsmaster/fs/user"
)

func Handle(path string, u *user.User, body string) []byte {
	fn := FileName(body)
	result := common.MessageInFS{}
	switch path {

	case "/fs/mkdir":
		succ, _, err := fn.MakeDir(u)
		result.Success = succ
		if err != nil {
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

func Load(file string, size int64, user *user.User) *common.GFSWriter {
	blocks := utils.Ceil(float64(size) / float64(common.BlockSize))
	logger.Info(blocks, size)
	fls := make([]*common.FileLocation, 0, blocks)
	for i := 0; i < blocks; i++ {
		fls = append(fls, &common.FileLocation{Main: "http://localhost:8087"})
	}
	fn := FileName(file)
	tFile := fn.Find(user)
	tFile.Size = size
	tFile.Location = fls
	return &common.GFSWriter{TargetFile: file, TargetNode: fls, MaxSize: size}
}

func Get(file string, user *user.User) (*common.GFSReader, error) {
	fn := FileName(file)
	tFile := fn.Find(user)
	if tFile != nil && !tFile.Unvisiable {
		return &common.GFSReader{TargetFile: file, TargetNode: tFile.Location, MaxSize: tFile.Size}, nil
	} else {
		return nil, errors.New("请等待")
	}
}
