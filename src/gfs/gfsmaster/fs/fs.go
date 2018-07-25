package fs

import (
	"gfs/common"
	"gfs/common/gfs"
	"gfs/common/glog"
	"gfs/common/gutils"
	"gfs/gfsmaster/common"
	"gfs/gfsmaster/wal"
	"os"
	"strconv"
	"sync"
	"time"
)

type GFS struct {
	sync.Mutex
	FileSystem *GFileSystem
	Users      map[string]gfs.User
	SerialId   int
}

var logger = glog.GetLogger("gfs/gfsmaster/fs")

var users = make(map[string]gfs.User)
var MyGFS = GFS{SerialId: 1, FileSystem: myfs, Users: users}

//从gonf中读取配置
func InitSystem() {

}

func RestoreFromLocal() {
	storageConf := mcommon.GetPeerConf().Storage
	fullName := storageConf.Image + "/image_full"
	file, err := os.OpenFile(fullName, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic("fuck")
	}
	common.DecodeFromReader(&MyGFS, file)
	buf := make([]*wal.WAL, 50, 50)
	cTId := MyGFS.FileSystem.current()
	for {
		ind, err := wal.ReadWAL(cTId, buf)
		if err != nil {
			panic("fuck")
		}
		if ind == 0 {
			break
		}
		for _, ww := range buf[0:ind] {
			MyGFS.FileSystem.replayWAL(ww, MyGFS.Users[ww.User])
		}
	}
}

func RestoreFromRemote() {

}

func Start() {

}
func storeImage() {
	storageConf := mcommon.GetPeerConf().Storage
	fileName := storageConf.Image + string(os.PathSeparator) + "image_" + gutils.GetNowStringSimple()
	interval := storageConf.Interval
	file, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	common.EncodeToWriter(MyGFS, file)
	file.Close()
	fullName := storageConf.Image + string(os.PathSeparator) + "image_full"
	if _, err := os.Stat(fullName); err == nil {
		os.Remove(fullName)
	}
	os.Rename(fileName, fullName)
	time.AfterFunc(time.Duration(interval)*time.Second, func() { storeImage() })
}

func MkDir(path string, recurisive bool, user gfs.User) (*gfs.File, error) {
	MyGFS.Lock()
	defer MyGFS.Unlock()
	wal.WriteWAL(wal.Create(MyGFS.FileSystem.next(), path, "", "r:"+strconv.FormatBool(recurisive), "", user.GetName(), wal.OP_CREATE_DIR))
	return MyGFS.FileSystem.MkDir(path, recurisive, user)
}
func Touch(path string, user gfs.User) (*gfs.File, error) {
	MyGFS.Lock()
	defer MyGFS.Unlock()
	wal.WriteWAL(wal.Create(MyGFS.FileSystem.next(), path, "", "", "", user.GetName(), wal.OP_CREATE_FILE))
	return MyGFS.FileSystem.Touch(path, user)
}
func Exists(path string, user gfs.User) (bool, error) {
	return MyGFS.FileSystem.Exists(path, user)
}
func List(path string, user gfs.User) ([]*gfs.File, error) {
	return MyGFS.FileSystem.List(path, user)
}
func GetFileInfo(path string, user gfs.User) (*gfs.File, error) {
	return MyGFS.FileSystem.GetFileInfo(path, user)
}
func Remove(path string, recurisive bool, user gfs.User) (*gfs.File, error) {
	MyGFS.Lock()
	defer MyGFS.Unlock()
	wal.WriteWAL(wal.Create(MyGFS.FileSystem.next(), path, "", "r:"+strconv.FormatBool(recurisive), "", user.GetName(), wal.OP_MV_FILE))
	return MyGFS.FileSystem.Remove(path, recurisive, user)
}
