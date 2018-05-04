package server

import (
	"fmt"
	"gfs/common"
	"gfs/gfsnode/data"
	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

var logger = logging.MustGetLogger("gfs/gfsnode/server")

func Cmd() *cobra.Command {
	var conf string
	var cmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			if c, err := common.GetConf(conf); err == nil {
				server := Server(*c)
				server.start()
			}
		},
	}
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "配置文件位置")
	return cmd
}

type Server common.Conf

func (server *Server) start() {
	node := server.Node
	http.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		fileName := req.FormValue("file")
		fileBlock := req.FormValue("block")

		if b, err := ioutil.ReadAll(req.Body); err == nil {
			d := &data.Data{File: fileName, Block: fileBlock, Data: b}
			d.Store(&(server.Node))
		}
	})
	http.ListenAndServe(":"+node.AdvisePort, nil)
	ioutil.ReadAll(r)
}
