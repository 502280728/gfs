// 用于datanode本地的文件存储
// 主要是保存本地存储的文件名，对应的blockid，相应的时间，本地存储位置等等
package data

import (
	"container/list"
	"fmt"
	"gfs/common"
	http1 "gfs/common/http"
	"gfs/common/utils"
	"os"
	"strconv"
	"sync"
	"time"
)

//代表一个文件块（block）
type Data struct {
	File       string    //原始的存在master上面的文件名
	Block      int       //第几段
	Path       string    //本地文件
	reported   bool      //是否报告给master
	modifyTime time.Time //该block被最后修改的时间,给master报告的时候，如果modifyTime与当前时间的差小于5秒，那么等待下一次报告
	creatTime  time.Time //该文件块的创建时间
}

//代表datanode整个本地缓存
type FileStore struct {
	sync.RWMutex
	Cache          map[string]*list.Element
	list           *list.List        //双端队列，最新写入的放到最前面
	noticeInterval int               //每隔多久通知master，单位是秒
	master         string            //对应的通知的master url，格式是http://host:port
	localDir       string            //该缓存在本地的持久化存储路径
	dataDir        string            //block在本地的存储位置
	config         map[string]string //其他的配置项
}

//配置FileStore
//至少需要master的配置
func (fs *FileStore) Config(config map[string]string) {
	if mh, found := config["master"]; found {
		var p string
		if mp, f := config["masterPort"]; f {
			p = mp
		} else {
			p = "8081"
		}
		fs.master = mh + ":" + p
	} else {
		panic("can not find master config") //不配置master根本无法工作
	}

	if ni, found := config["noticeInterval"]; found {
		tmp, _ := strconv.Atoi(ni)
		fs.noticeInterval = tmp
	} else {
		fs.noticeInterval = 10
	}

	if ld, found := config["localDir"]; found {
		fs.localDir = ld
	} else {
		fs.localDir = os.TempDir() + "/gfsnode/cache"
	}

	if dd, found := config["dataDir"]; found {
		fs.dataDir = dd
	} else {
		fs.dataDir = os.TempDir() + "/gfsnode/data"
	}
	fs.config = config
	fs.Cache = make(map[string]*list.Element)
	fs.list = list.New()
	logger.Infof("datanode config, master:  %s", fs.master)
	logger.Infof("datanode config, localDir:  %s", fs.localDir)
	logger.Infof("datanode config, dataDir:  %s", fs.dataDir)
	logger.Infof("datanode config, noticeInterval:  %s", fs.noticeInterval)

}

//启动整个FileStore
func (fs *FileStore) Init() {
	utils.CreateDirIfNotExists(fs.localDir)
	utils.CreateDirIfNotExists(fs.dataDir)
	fs.NoticeMaster()
	fs.Persistent()
}

func (fs *FileStore) persistent() {
	fs.Lock()
	bb := common.EncodeToBytes(*fs)
	if file, err := os.OpenFile(fs.localDir+"/blocks.data", os.O_TRUNC, os.ModePerm); err == nil {
		file.Write(bb)
		file.Close()
	} else {
		logger.Errorf("error when persist data to %s with error %s", fs.localDir, err.Error())
	}
	fs.Unlock()
}
func (fs *FileStore) Persistent() {
	logger.Infof("persist data to %s with interval %s seconds", fs.localDir, 20)
	fs.persistent()
	time.AfterFunc(20*time.Second, func() {
		fs.Persistent()
	})
}

func (fs *FileStore) CheckIfExists(key string) (*Data, bool) {
	fs.RLock()
	defer fs.RUnlock()
	if d, found := fs.Cache[key]; found {
		return d.Value.(*Data), true
	} else {
		return nil, false
	}
}

func (fs *FileStore) Set(key string, data *Data) {
	fs.Lock()
	defer fs.Unlock()
	ele := fs.list.PushFront(data)
	fs.Cache[key] = ele
}

func (fs *FileStore) noticeMaster() {
	fs.Lock()
	notices := []*Data{}
	for {
		if ele := fs.list.Front(); ele != nil {
			data := ele.Value.(*Data)
			if !data.reported {
				if time.Now().Unix()-data.modifyTime.Unix() >= 5 {
					notices = append(notices, data)
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	var req = http1.GFSRequest{}
	var res common.ACK
	if err := req.PostObj(fs.master+"/node", notices, &res); err == nil {
		if res {
			for _, n := range notices {
				n.reported = true
			}
		}
	} else {
		logger.Errorf("error when notice %s with error %s", fs.master, err.Error())
	}
	fs.Unlock()
}

func (fs *FileStore) NoticeMaster() {
	logger.Infof("noticing master %s ever %d seconds", fs.master, fs.noticeInterval)
	fs.noticeMaster()
	time.AfterFunc(time.Duration(int64(fs.noticeInterval)*int64(time.Second)), func() {
		fs.NoticeMaster()
	})
}
