package data

import (
	"fmt"
	"gfs/common"
	"github.com/satori/uuid"
	"io/ioutil"
	"log"
	"os"
)

type Data struct {
	File  string //原始的存在master上面的文件名
	Block string //第几段
	Path  string //本地文件
	Data  []byte // 数据
}

//
func (data *Data) Store(conf *common.DataNode) {

	if tmp, err := uuid.NewV4(); err == nil {
		data.Path = fmt.Sprintf("%s", tmp)
		log.Printf("store file in %s", data.Path)
		ioutil.WriteFile(conf.DataDir+"/"+data.Path, data.Data, os.ModeAppend)
	}
}
