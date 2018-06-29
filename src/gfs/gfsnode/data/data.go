package data

import (
	"fmt"
	"gfs/common"
	"gfs/common/utils"
	"log"
	"os"
	"strconv"
	"time"
)

func (data *Data) Store(conf *common.DataNode, bb []byte) {
	if d, found := LocalStore.CheckIfExists(data.File + strconv.Itoa(data.Block)); found {
		storeFile(d.Path, bb)
		d.modifyTime = time.Now()
	} else {
		if tmp, err := utils.UUID(); err == nil {
			data.Path = fmt.Sprintf("%s/%s", conf.DataDir, tmp)
			log.Printf("store file in %s", data.Path)
			storeFile(data.Path, bb)
			LocalStore.Set(data.File+strconv.Itoa(data.Block), &Data{File: data.File, Block: data.Block, Path: data.Path, reported: false, modifyTime: time.Now()})
		}
	}
}

func storeFile(file string, bb []byte) {
	if _, err := os.Stat(file); err != nil {

		os.Create(file)
	}
	ff, _ := os.OpenFile(file, os.O_APPEND, os.ModeAppend)
	defer ff.Close()
	ff.Write(bb)
}

func (data *Data) Retrieve(conf *common.DataNode, fgc *common.FileBlockChip) {
	if d, found := LocalStore.CheckIfExists(data.File + strconv.Itoa(fgc.Block)); found {
		logger.Infof("begin retrieve block %d of file %s in local store ", fgc.Block, data.File)
		retrieveFromBlockFile(d.Path, fgc)
		logger.Infof("end retrieve block %d of file %s in local store ", fgc.Block, data.File)
	} else {
		logger.Errorf("can not find block %d of file %s in local store", fgc.Block, data.File)
	}
}

func retrieveFromBlockFile(f string, fgc *common.FileBlockChip) {
	file, err := os.Open(f)
	defer file.Close()
	if err == nil {
		file.Seek(fgc.Offset, 0)
		fgc.Data = make([]byte, fgc.Limit, fgc.Limit)
		size, _ := file.Read(fgc.Data)
		fgc.Data = fgc.Data[0:size]
	} else {
		logger.Error(err.Error())
	}
}
