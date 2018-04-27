package server

import (
	"gfs/gfsmaster/fs"
	"gfs/gfsmaster/fs/user"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Cmd() *cobra.Command {
	var conf string
	var cmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			initFileSystem()
			createServer()
		},
	}
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "配置文件位置")
	return cmd
}

func initFileSystem() {
	fs.RootStorePath = "D:/temp/fs1.binary"
	fs.RecoverFromStore()
}

type Handler func(http.ResponseWriter, *http.Request)

func (handler Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler(w, req)
}

func createServer() {
	http.ListenAndServe(":8080", createHandler())
}

func createHandler() Handler {
	handler := Handler(func(w http.ResponseWriter, req *http.Request) {
		uri, _ := url.Parse(req.RequestURI)
		if strings.HasPrefix(uri.Path, "/fs") {
			userName := uri.Query().Get("user")
			fs.CreateHandler(uri.Path, &user.User{Name: userName}, parseBody(req.Body))
		}

	})
	return handler
}

func parseBody(body io.ReadCloser) string {
	b, _ := ioutil.ReadAll(body)
	return string(b)
}
