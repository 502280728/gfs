package server

import (
	"gfs/common"
	"gfs/gfsnode/data"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

func Cmd() *cobra.Command {
	var conf string
	var cmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			if c, err := getConf(conf); err == nil {
				server := Server(*c)
				server.start()
			}
		},
	}
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "配置文件位置")
	return cmd
}

type Server common.Conf

func getConf(file string) (*common.Conf, error) {
	if b, err := ioutil.ReadFile(file); err == nil {
		conf := &common.Conf{}
		if err := yaml.Unmarshal(b, conf); err == nil {
			return conf, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (server *Server) start() {
	node := server.Node
	http.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		fileName := req.FormValue("file")
		fileBlock := req.FormValue("block")

		if b, err := ioutil.ReadAll(req.Body); err == nil {
			//fmt.Fprint(w, string(b)+"asdfsdf")
			d := &data.Data{File: fileName, Block: fileBlock, Data: b}
			d.Store(&(server.Node))
		}
	})
	http.ListenAndServe(":"+node.Port, nil)
}
