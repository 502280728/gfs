package wal

import (
	"bytes"
	"gfs/common/glog"
	"gfs/gfsmaster/common"
	"strconv"
	"strings"
)

var logger = glog.GetLogger("gfs/gfsmaster/wal")

type WAL struct {
	ID         string //事务id
	Op         Operate
	SourceFile string //一个操作对应的源文件
	TargetFile string //一个操作对应的目标文件，如果存在的话；比如移动文件时的目标文件
	Params     string //命令的参数,比如删除文件夹时是否启用了递归删除
	Extra      string //其他的信息,具体格式每个Op可能都不一样，比如在ALLOCATE_FILE时，要记录下分配了哪些block
	User       string
}

func (w *WAL) ToBytes() []byte {
	var bb bytes.Buffer
	bb.WriteString(w.ID)
	bb.WriteString(" ")
	bb.WriteString(strconv.Itoa(int(w.Op)))
	bb.WriteString(" ")
	bb.WriteString(w.SourceFile)
	bb.WriteString(" ")
	bb.WriteString(w.TargetFile)
	bb.WriteString(" ")
	bb.WriteString(w.Params)
	bb.WriteString(" ")
	bb.WriteString(w.Extra)
	bb.WriteString(" ")
	bb.WriteString(w.User)
	bb.WriteString(" ")
	bb.WriteByte('\r')
	bb.WriteByte('\n')
	return bb.Bytes()
}

func (res *WAL) Parse(line string) {
	ss := strings.Split(line, " ")
	res.ID = ss[0]
	oo, _ := strconv.Atoi(ss[1])
	res.Op = Operate(oo)
	res.SourceFile = ss[2]
	res.TargetFile = ss[3]
	res.Params = ss[4]
	res.Extra = ss[5]
	res.User = ss[6]
}

//wal的类型,也就是操作类型
type Operate int

var walChan = make(chan *WAL)
var walBuf = make([]*WAL, 0, 1000)
var walConf *mcommon.WALConf
var walWritor WALWritor
var walReader WALReader

const (
	OP_RM_FILE       Operate = 1 << iota //删除文件或者文件夹
	OP_CREATE_FILE                       //新建文件
	OP_CREATE_DIR                        //新建文件夹
	OP_MV_FILE                           //移动文件
	OP_MV_DIR                            //移动文件夹
	OP_ALLOCATE_FILE                     //为某个文件分配block
)

func Create(id, sf, tf, params, extra, user string, op Operate) *WAL {

	return &WAL{
		ID:         id,
		Op:         op,
		SourceFile: sf,
		TargetFile: tf,
		Params:     params,
		Extra:      extra,
		User:       user,
	}
}

func InitWAL() {
	walConf = mcommon.GetPeerConf().WAL
	walReader = &WALFileReader{}
}

func Start() {
	logger.Info("start goroutine for wal")
	if walConf.Fsync == "everysec" {
		walWritor = &WALBufWritor{buf: make([]*WAL, 0, 100)}
	} else if walConf.Fsync == "always" {
		walWritor = &WALChanWritor{buf: make(chan *WAL)}
	}
	go walWritor.Flush()
}

func WriteWAL(wal *WAL) {
	walWritor.Write(wal)
}

func ReadWAL(tid string, buf []*WAL) (int, error) {
	return walReader.Read(tid, buf)
}
