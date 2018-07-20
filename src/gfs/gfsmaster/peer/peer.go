// peer
package peer

import (
	"gfs/common"
	"gfs/common/ghttp"
	"gfs/common/glog"
	"gfs/gfsmaster/common"
	"gfs/gfsmaster/fs"
)

var logger = glog.GetLogger("gfs/gfsmaster/peer")

type PeerType uint8
type PeerStatus int16

//master的节点类型
const (
	Leader   PeerType = 1 << iota //leader类型
	Follower                      //follower类型
)

const (
	Empty     PeerStatus = 1 << iota //初始状态，相当于什么都没有的状态，没有实际的意义，仅是为了让状态更加全
	Initing                          //读取配置中
	Inited                           //读取配置完成，这是个瞬时态
	Restoring                        //数据恢复中,master正在恢复数据，follower正在获取leader的数据，此时系统不能对外提供服务
	Restored                         //数据恢复完成，这基本也是个瞬时态
	Working                          //工作中，可以对外提供服务
	Stopping                         //关闭中，在做清理工作，如将WAL，gdb刷到硬盘

)

type Peer struct {
	Type   PeerType
	Conf   *mcommon.PeerConf
	Status PeerStatus
}

//节点对象
var peer Peer

//
func InitPeer(confFile ...string) {
	peer.loadConf(confFile...)
	fs.InitSystem()
}

func GetPeer() *Peer {
	return &peer
}

func StartPeer() {
	peer.start()
}

// 核心方法,用于启动master节点
func (peer *Peer) start() {
	peer.startServer()       //在另一个goroutine上启动,能够接受外界的请求了
	peer.restoreFileSystem() //
	peer.startFileSystem()   //会启动filesystem需要的一些goroutine,包括定时写image，何时写WAL等等
}

//加载配置
func (peer *Peer) loadConf(conFile ...string) {
	peer.statusTo(Initing)
	mcommon.LoadConf(conFile...)
	peer.Conf = mcommon.GetPeerConf()
	peer.checkIfLeader()
	peer.statusTo(Inited)
	logger.Info("master finished reading confs")
}

//恢复数据
func (peer *Peer) restoreFileSystem() {
	peer.statusTo(Restoring)
	//TODO
	fs.RestoreFromLocal()
	if peer.Type == Follower {
		fs.RestoreFromRemote()
	}
	peer.statusTo(Restored)
}

//启动服务
func (peer *Peer) startServer() {
	//TODO
	peer.statusTo(Working)
}

func (peer *Peer) startFileSystem() {
	fs.Start()
}

//改变master节点的状态
func (peer *Peer) statusTo(status PeerStatus) {
	peer.Status = status
}

func (peer *Peer) checkIfLeader() {
	if peer.Conf.Address.FriendIP == "" {
		peer.Type = Leader
	} else {
		var http ghttp.GFSRequest
		var res common.ACK
		address := peer.Conf.Address
		err := http.PostObj(address.Protocol+"://"+address.FriendIP+":"+address.M2MPort, common.ACK(true), res)
		if err == nil && bool(res) {
			peer.Type = Follower
		} else {
			peer.Type = Leader
		}
	}
}
