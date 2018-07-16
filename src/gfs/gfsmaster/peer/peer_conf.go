// peer_conf
package peer

import (
	"gfs/common/gconf"
	"gfs/common/gutils"

	"gfs/common/glog"
)

var logger = glog.GetLogger("gfs/gfsmaster/peer")

const (
	GFS_MASTER_IP        = "gfs.master.ip"
	GFS_MASTER_FRIEND_ID = "gfs.master.friend.ip"
	GFS_MASTER_M2M_PORT  = "gfs.master.m2m.port"
	GFS_MASTER_M2C_PORT  = "gfs.master.m2c.port"
	GFS_MASTER_M2D_PORT  = "gfs.master.m2d.port"

	GFS_MASTER_LOG        = "gfs.master.log"
	GFS_MASTER_LOG_FILE   = "gfs.master.log.file"
	GFS_MASTER_LOG_LEVEL  = "gfs.master.log.level"
	GFS_MASTER_LOG_FORMAT = "gfs.master.log.format"

	GFS_MASTER_WAL       = "gfs.master.wal"
	GFS_MASTER_WAL_FILE  = "gfs.master.wal.file"
	GFS_MASTER_WAL_FSYNC = "gfs.master.wal.fsync"

	GFS_MASTER_FILE          = "gfs.master.file"
	GFS_MASTER_FILE_INTERVAL = "gfs.master.interval"
)

type PeerConf struct {
	Storage *StorageConf
	Log     *LogConf
	WAL     *WALConf
	Address *AddressConf
}

//master本地文件的存储配置项
type StorageConf struct {
	File     string //存储位置
	Interval int    //每隔多少秒存储一次
}

//网络地址与端口的配置项
type AddressConf struct {
	IP       string //当前机器的ip
	FriendIP string //双机master方案中另一个master地址
	Protocol string //通信使用的协议，默认是http
	M2MPort  string //master与master之间通信的端口
	M2CPort  string //master与client之间通信的端口
	M2DPort  string //master与datanode之间通信的端口
}

//日志配置项
type LogConf struct {
	Log    string //on or off
	File   string //文件位置
	Level  string //日志级别
	Format string
}

//WAL配置项
type WALConf struct {
	WAL   string //on or off
	File  string //文件位置
	Fsync string //always or everysec
}

func (conf *PeerConf) LoadConfs(confFile ...string) {
	gconf.Conf.LoadConfs(confFile...)
	conf.Log = &LogConf{}
	conf.Log.Load()
	conf.Storage = &StorageConf{}
	conf.Storage.Load()
	conf.Address = &AddressConf{}
	conf.Address.Load()
	conf.WAL = &WALConf{}
	conf.WAL.Load()
}

func (conf *StorageConf) Load() {
	gconf.Conf.GetOrDefault(GFS_MASTER_FILE, &conf.File, "temp/master.gob")
	gconf.Conf.GetOrDefault(GFS_MASTER_FILE_INTERVAL, &conf.Interval, 60)
	logger.Debugf("storage conf: file -> %s", conf.File)
	logger.Debugf("storage conf: interval -> %d", conf.Interval)
}

func (conf *LogConf) Load() {
	gconf.Conf.GetOrDefault(GFS_MASTER_LOG, &conf.Log, "on")
	gconf.Conf.GetOrDefault(GFS_MASTER_LOG_FILE, &conf.File, "temp/logs.log")
	gconf.Conf.GetOrDefault(GFS_MASTER_LOG_LEVEL, &conf.Level, "info")
	gconf.Conf.GetOrDefault(GFS_MASTER_LOG_FORMAT, &conf.Format, "%{time:2006-01-02 15:04:05} %{module} %{level} %{message}")
	glog.File = conf.File
	glog.Level = conf.Level
	glog.Format = conf.Format
	glog.LogOn = conf.Log == "on"
	glog.Init()
	logger.Debugf("log conf: file -> %s", conf.File)
	logger.Debugf("log conf: level -> %s", conf.Level)
	logger.Debugf("log conf: format -> %s", conf.Format)
}

func (conf *AddressConf) Load() {
	gconf.Conf.GetOrDefault(GFS_MASTER_IP, &conf.IP, gutils.GetLocalIP())
	gconf.Conf.Get(GFS_MASTER_FRIEND_ID, &conf.FriendIP)
	gconf.Conf.GetOrDefault(GFS_MASTER_M2M_PORT, &conf.M2MPort, "40001")
	gconf.Conf.GetOrDefault(GFS_MASTER_M2C_PORT, &conf.M2CPort, "40002")
	gconf.Conf.GetOrDefault(GFS_MASTER_M2D_PORT, &conf.M2DPort, "40003")
	logger.Debugf("master conf: ip -> %s", conf.IP)
	logger.Debugf("master conf: friend ip -> %s", conf.FriendIP)
	logger.Debugf("master conf: m2mport -> %s", conf.M2MPort)
	logger.Debugf("master conf: m2cport -> %s", conf.M2CPort)
	logger.Debugf("master conf: m2dport -> %s", conf.M2DPort)
}
func (conf *WALConf) Load() {
	gconf.Conf.GetOrDefault(GFS_MASTER_WAL, &conf.WAL, "on")
	gconf.Conf.GetOrDefault(GFS_MASTER_WAL_FILE, &conf.File, "temp/wal.wal")
	gconf.Conf.GetOrDefault(GFS_MASTER_WAL_FSYNC, &conf.Fsync, "everysec")
	logger.Debugf("wal conf: wal -> %s", conf.WAL)
	logger.Debugf("wal conf: file -> %s", conf.File)
	logger.Debugf("wal conf: fsync -> %s", conf.Fsync)
}
