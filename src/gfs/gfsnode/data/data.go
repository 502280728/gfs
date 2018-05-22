package data

import (
	"fmt"
	"gfs/common"
	logging "github.com/op/go-logging"
	"github.com/satori/uuid"
	"log"
	"os"
	"strconv"
)

var logger = logging.MustGetLogger("gfs/gfsnode/data")

type Data struct {
	File  string //原始的存在master上面的文件名
	Block int    //第几段
	Path  string //本地文件
}

type FileStore struct {
	Cache map[string]*Data
}

var LocalStore FileStore

func init() {
	LocalStore = FileStore{Cache: map[string]*Data{}}
	//WZF this is only for test
	LocalStore.Cache["/wzf/a.txt0"] = &Data{File: "/wzf/a.txt", Block: 0, Path: "d:/temp/a.txt"}
}

//
func (data *Data) Store(conf *common.DataNode) {
	if tmp, err := uuid.NewV4(); err == nil {
		data.Path = fmt.Sprintf("%s", tmp)
		log.Printf("store file in %s", data.Path)
		//ioutil.WriteFile(conf.DataDir+"/"+data.Path, data.Data, os.ModeAppend)
	}
}

func (data *Data) Retrieve(conf *common.DataNode, fgc *common.FileBlockChip) {
	if d, found := LocalStore.Cache[data.File+strconv.Itoa(fgc.Block)]; found {
		logger.Infof("begin retrieve block %s of file %s in local store ", fgc.Block, data.File)
		retrieveFromBlockFile(d.Path, fgc)
		logger.Infof("end retrieve block %s of file %s in local store ", fgc.Block, data.File)
	} else {
		logger.Errorf("can not find block %s of file %s in local store", fgc.Block, data.File)
	}
}

func retrieveFromBlockFile(f string, fgc *common.FileBlockChip) {
	file, err := os.Open(f)
	defer file.Close()
	if err == nil {
		file.Seek(fgc.Offset, 0)
		fgc.Data = make([]byte, fgc.Limit, fgc.Limit)
		file.Read(fgc.Data)
	} else {
		logger.Error(err.Error())
	}
}
