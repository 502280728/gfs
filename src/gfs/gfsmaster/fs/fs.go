package fs

import (
	"gfs/common"
	"gfs/common/gfs"
	"gfs/common/gio"
	"gfs/common/glog"
	"gfs/common/gutils"
	"gfs/gfsmaster/common"
	"os"
	"time"
)

type GFS struct {
	FileSystem *GFileSystem
	Users      []gfs.User
	SerialId   int
}

var wal = make(chan string, 400)

var logger = glog.GetLogger("gfs/gfsmaster/fs")

var users = make([]gfs.User, 0, 10)
var MyGFS = GFS{SerialId: 1, FileSystem: myfs, Users: users}

//从gonf中读取配置
func InitSystem() {

}

func RestoreFromLocal() {
	storageConf := mcommon.GetPeerConf().Storage
	fullName := storageConf.Image + "/image_full"
	if file, err := os.OpenFile(fullName, os.O_RDONLY, os.ModeAppend); err == nil {
		common.DecodeFromReader(&MyGFS, file)
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
func storeWAL() {
}

func Create(path string, user gfs.User) (gio.WriteCloser, error) {
	return MyGFS.FileSystem.Create(path, user)
}
func Open(path string, user gfs.User) (gio.ReadCloser, error) {
	return MyGFS.FileSystem.Open(path, user)
}
func MkDir(path string, recurisive bool, user gfs.User) (*gfs.File, error) {
	return MyGFS.FileSystem.MkDir(path, recurisive, user)
}
func Touch(path string, user gfs.User) (*gfs.File, error) {
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
	return MyGFS.FileSystem.Remove(path, recurisive, user)
}
