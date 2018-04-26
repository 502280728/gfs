package common

import ()

type DataNode struct {
	DataDir      string   //`datanode上文件的存储位置`
	BlockSize    string   //`每个文件块的大小`
	InfoInterval int      //`间隔n秒向master报告自己的情况`
	Masters      []string //`master的URI`
	Port         string   //`开放的端口，用于接收master发过来的数据`
}

type MasterNode struct {
	DefaultDir string //master节点存放文件位置
	DefaultFs  string //master节点的url，默认是localhost:9090
}

type Conf struct {
	Node   DataNode
	Master MasterNode
}
