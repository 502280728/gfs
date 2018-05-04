package node

//master的datanode节点管理
import (
	"bytes"
	"encoding/gob"
	"gfs/common"
	"gfs/gfsmaster/fs"
	"gfs/gfsmaster/fs/user"
	logging "github.com/op/go-logging"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var logger = logging.MustGetLogger("gfs/gfsmaster/node")

//node到master中报备信息，key是node的adviseaddress，value是报备的时间戳
var DataNodeCache = make(map[string]time.Time)

//key是文件名，value是一个切片，假设一个文件被分成了10块，那么这个切片长度是10，capacity也是10
//切片的第一个位置放置第一个block的存放位置，如果有多个用分号分隔;第二个位置存放第二个block的存放位置，以此类推
//该map的value必须在master分配地址时初始化，使value的切片长度为该文件包含的block的个数
var FileLocation = make(map[string][]string)

type File struct {
	name          string //文件名
	blockid       int    //blockid
	adviseaddress string //datanode的host+port
}

var channel = make(chan File)

func DoSomework() {
	go func(chan []File) {
		tmp := <-channel
		if FileLocation[tmp.name] != nil {
			if len(FileLocation[tmp.name][tmp.blockid]) != 0 {
				FileLocation[tmp.name][tmp.blockid] = FileLocation[tmp.name][tmp.blockid] + ";" + adviseaddress
			} else {
				FileLocation[tmp.name][tmp.blockid] = adviseaddress
			}
		}
	}(channel)
}

func HandleNodeRequest(adviseAddress string, bb []byte) []byte {
	logger.Infof("receive info from %s", adviseAddress)
	DataNodeCache[adviseAddress] = time.Now()
	var dnim common.DataNodeIntervalMessage
	dec := gob.NewDecoder(bb)
	dec.Decode(&dnim)
	files := dnim.Files
	for _, file := range files {
		ss := strings.Split(file, ":")
		bid, _ := strconv.Atoi(ss[1])
		channel <- File{adviseaddress: adviseAddress, blockid: bid, name: ss[0]}
	}
	var res bytes.Buffer
	rs := common.ACK(true)
	enc := gob.NewEncoder(res)
	enc.Encode(rs)
	return res.Bytes()
}

func HandleClientRequest(filename string, blocks int, u *user.User) []byte {
	//TODO 在文件树中新建这个文件
	fn := fs.FileName(filename)
	fn.Touch(u)
	var mtcm common.MasterToClientMessage
	var tmp []string
	for k, v := range DataNodeCache {
		if time.Now().Unix()-v.Unix() <= 60 {
			tmp = append(tmp, k)
		}
	}
	for i := 0; i < blocksize; i++ {
		var tmpstr = make([]string, 3, 3)
		for j := 0; j < 3; j++ {
			tmpstr[j] = tmp[rand.Intn(len(tmp))]
		}
		mtcm.Nodes = append(mtcm.Nodes, strings.Join(tmpstr, ";"))
	}
	FileLocation[filename] = make([]string, blocksize, blocksize)
	var res bytes.Buffer
	enc := gob.NewEncoder(mtcm)
	enc.Encode(rs)
	return res.Bytes()
}
