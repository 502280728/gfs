package common

import (
	"errors"
	logging "github.com/op/go-logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var logger = logging.MustGetLogger("gfs/common")

//在datanode上每个文件块的大小，单位是字节
var BlockSize int64 = 1024

type DataNode struct {
	DataDir      string   //datanode上文件的存储位置
	InfoInterval int      //间隔n秒向master报告自己的情况
	Masters      []string //master的URI
	AdvisePort   string   //开放的端口，用于接收master发过来的数据
	AdviseHost   string   //建议master连接使用的ip或者地址
}

type MasterNode struct {
	DefaultDir string //master节点存放文件位置
	DefaultFs  string //master节点的url，默认是localhost:9090
	BlockSize  string //每个文件块的大小
}

type Conf struct {
	Node   DataNode
	Master MasterNode
}

func GetConf(file string) (*Conf, error) {
	if res, err := ioutil.ReadFile(filename); err == nil {
		var conf = &Conf{}
		err := yaml.Unmarshal([]byte(res), conf)
		if err != nil {
			return nil, err
		} else {
			return conf, nil
		}
	} else {
		return nil, err
	}
}
