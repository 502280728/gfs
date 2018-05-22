package server

import (
	"gfs/common"
	"gfs/gfsnode/data"
	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strconv"
)

var logger = logging.MustGetLogger("gfs/gfsnode/server")

type Server common.Conf

func Cmd() *cobra.Command {
	var conf string
	var cmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("start server,using conf file %s", conf)
			if c, err := common.GetConf(conf); err == nil {
				server := Server(*c)
				server.start()
			} else {
				logger.Errorf("errors occurs when reading file %s", conf)
				logger.Error(err)
			}
		},
	}
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "配置文件位置")
	return cmd
}

func (server *Server) start() {
	node := server.Node
	logger.Infof("start server in port : %s", node.AdvisePort)
	http.HandleFunc("/data/in", createDataInHandler(server))
	http.HandleFunc("/data/out", createDataOutHandler(server))
	http.ListenAndServe(":"+node.AdvisePort, nil)

}

//往datanode写入数据
func createDataInHandler(svrConf *Server) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fileName := req.FormValue("file")
		fileBlock, _ := strconv.Atoi(req.FormValue("block"))
		if _, err := ioutil.ReadAll(req.Body); err == nil {
			d := &data.Data{File: fileName, Block: fileBlock}
			d.Store(&(svrConf.Node))
		}
	}
}

//从datanode中读取数据
func createDataOutHandler(svrConf *Server) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if bb, err := ioutil.ReadAll(req.Body); err == nil {
			var fgc = &common.FileBlockChip{}
			fgc.Decode(bb)
			logger.Infof("received param %s", fgc)
			d := &data.Data{File: fgc.FileName}
			d.Retrieve(&(svrConf.Node), fgc)
			w.Write(fgc.Data)
		}
	}
}
